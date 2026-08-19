package app

import (
	"context"
	"log/slog"
	"net"

	"github.com/diyorbeknematov/minitwitter/gen/go/tweet"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/repository/postgres"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/service"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/traceid"
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
	db, err := postgres.ConnectionDB(cfg.DB)
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to connect to database", err)
	}

	// Repository
	repo := repository.NewRepository(db)

	// Services
	svc := service.NewService(repo, cfg, logger)

	// gRPC Server
	grpcServer := grpc.NewServer(
		grpc.UnaryInterceptor(traceid.TraceServerInterceptor),
	)

	tweet.RegisterTweetServiceServer(grpcServer, svc.Tweet)

	// Listener
	listener, err := net.Listen("tcp", cfg.GRPC.Address())
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
		slog.Int("address", a.cfg.GRPC.Port),
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
		a.logger.Warn("graceful shotdown timeout, forcing stop")
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
