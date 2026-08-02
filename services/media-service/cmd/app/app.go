package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"

	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/repository/postgres"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/service"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/storage"
	"github.com/jmoiron/sqlx"
	"google.golang.org/grpc"
)

type App struct {
	logger *slog.Logger
	cfg    *config.Config

	db *sqlx.DB

	grpcServer *grpc.Server
	listener   net.Listener
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Database
	db, err := postgres.DBConnection(cfg.DB)
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to connect to database", err)
	}

	// Repository
	repo := repository.NewRepository(db)

	// Object Storage
	storage, err := storage.NewMinIO(cfg.MinIO)
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to initialize object storage", err)
	}

	// Services
	svc := service.NewMediaService(repo, storage, logger, cfg)

	// gRPC Server
	grpcServer := grpc.NewServer()

	media.RegisterMediaServiceServer(grpcServer, svc)

	// Listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPC.Port))
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to listen on gRPC port", err)
	}
	return &App{
		logger:     logger,
		cfg:        cfg,
		db:         db,
		grpcServer: grpcServer,
		listener:   listener,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info(
		"gRPC server started",
		slog.String("address", a.cfg.GRPC.Port),
	)

	return a.grpcServer.Serve(a.listener)
}

func (a *App) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		a.logger.Warn("graceful shutdown timeout, forcing stop")
		a.grpcServer.Stop()
	case <-done:
		a.logger.Info("gRPC server stopped")
	}

	if err := a.listener.Close(); err != nil {
		return apperror.Wrap("app", "Shutdown", "failed to close listener", err)
	}

	a.logger.Info("listener closed")

	if err := a.db.Close(); err != nil {
		return apperror.Wrap("app", "Shutdown", "failed to close database", err)

	}

	a.logger.Info("database connection closed")

	return nil
}
