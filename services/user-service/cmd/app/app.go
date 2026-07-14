package app

import (
	"context"
	"fmt"
	"net"

	"github.com/diyorbek/minitwitter/services/user-service/internal/config"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository/postgres"
	"github.com/diyorbek/minitwitter/services/user-service/internal/service"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/diyorbeknematov/minitwitter/gen/go/auth"
	"github.com/diyorbeknematov/minitwitter/gen/go/user"

	"log/slog"

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
	db, err := postgres.DBConnection(cfg)
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to connect to database", err)
	}

	// Repository
	repo := repository.NewRepository(db)

	// Services
	svc := service.NewService(repo, cfg, logger)

	// gRPC Server
	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, svc.Auth)
	user.RegisterUserServiceServer(grpcServer, svc.User)

	// Listener
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", cfg.GRPCPort))
	if err != nil {
		_ = db.Close()
		return nil, apperror.Wrap("app", "New", "failed to create tcp listener", err)
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
		slog.String("address", a.cfg.GRPCPort),
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

	if err := a.db.Close(); err != nil {
		return apperror.Wrap("app", "Shutdown", "failed to close database", err)
	}

	a.logger.Info("database connection closed")

	return nil
}
