package app

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"net/http"
	"time"

	"github.com/diyorbeknematov/minitwitter/gen/go/notification"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/repository/postgres"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/service"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/ws"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/pkg/apperror"
	"github.com/jmoiron/sqlx"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	logger *slog.Logger
	cfg    *config.Config

	db *sqlx.DB

	hub *ws.Hub

	grpcServer *grpc.Server
	grpcLn     net.Listener

	httpServer *http.Server
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// Database
	db, err := postgres.DBConnection(cfg.DB)
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to connect database", err)
	}

	// Repository
	repo := repository.NewRepo(db)

	// WebSocket Hub
	hub := ws.NewHub()

	// Service
	svc := service.NewService(repo, cfg, logger)

	// gRPC
	grpcServer := grpc.NewServer()

	notification.RegisterNotificationServiceServer(
		grpcServer,
		svc.Notification,
	)

	grpcLn, err := net.Listen(
		"tcp",
		fmt.Sprintf(":%d", cfg.Server.GRPCPort),
	)
	if err != nil {
		_ = db.Close()
		return nil, apperror.Wrap(
			"app",
			"New",
			"failed to create grpc listener",
			err,
		)
	}

	// HTTP (WebSocket)
	mux := http.NewServeMux()

	mux.Handle("/ws", ws.NewHandler(hub))

	httpServer := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.Server.HTTPPort),
		Handler: mux,
	}

	return &App{
		logger:     logger,
		cfg:        cfg,
		db:         db,
		hub:        hub,
		grpcServer: grpcServer,
		grpcLn:     grpcLn,
		httpServer: httpServer,
	}, nil
}

func (a *App) Run(ctx context.Context) error {
	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		a.logger.Info("gRPC server started", "port", a.cfg.Server.GRPCPort)

		if err := a.grpcServer.Serve(a.grpcLn); err != nil {
			return apperror.Wrap("app", "Run", "grpc server error", err)
		}
		return nil
	})

	g.Go(func() error {
		a.logger.Info("HTTP server started", "port", a.cfg.Server.HTTPPort)

		if err := a.httpServer.ListenAndServe(); err != nil &&
			err != http.ErrServerClosed {
			return apperror.Wrap("app", "Run", "http server error", err)
		}
		return nil
	})

	g.Go(func() error {
		<-ctx.Done()

		a.logger.Info("shutdown signal received, stopping app...")

		shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
		defer cancel()

		return a.Shutdown(shutdownCtx)
	})

	return g.Wait()
}

func (a *App) Shutdown(ctx context.Context) error {
	done := make(chan struct{})

	// gRPC
	go func() {
		a.grpcServer.GracefulStop()
		close(done)
	}()

	select {
	case <-ctx.Done():
		a.logger.Warn("graceful shudtown timeout, forcing stop")
		a.grpcServer.Stop()

	case <-done:
		a.logger.Info("gRPC server stopped")
	}

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return apperror.Wrap(
			"App",
			"Shutdown",
			"failed to shutdown http server",
			err,
		)
	}

	a.logger.Info("HTTP server stopped")

	// websocket hub
	a.hub.Close()

	// Database
	if err := a.db.Close(); err != nil {
		return apperror.Wrap(
			"app",
			"Shutdown",
			"failed to close database",
			err,
		)
	}

	a.logger.Info("database connection closed")

	if err := a.grpcLn.Close(); err != nil {
		a.logger.Warn("failed to close grpc listener", "error", err)
	}
	
	return nil
}
