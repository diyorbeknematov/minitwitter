package app

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/api"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/api/handler"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/service"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
)

const shutdownTimeout = 5 * time.Second

type App struct {
	logger *slog.Logger
	cfg    *config.Config

	grpcClient *grpcclient.Client

	httpServer *http.Server
}

func New(cfg *config.Config, logger *slog.Logger) (*App, error) {
	// gRPC clients
	grpcClient, err := grpcclient.New(cfg)
	if err != nil {
		return nil, apperror.Wrap("app", "New", "failed to create clients", err)
	}

	services := service.New(grpcClient)

	h := handler.New(services, logger, cfg.NotificationWSURL)

	router := api.New(h)

	httpServer := &http.Server{
		Addr:    cfg.Server.Address(),
		Handler: router,
	}

	return &App{
		logger:     logger,
		cfg:        cfg,
		grpcClient: grpcClient,
		httpServer: httpServer,
	}, nil
}

func (a *App) Run() error {
	a.logger.Info(
		"HTTP server started",
		"addr", a.httpServer.Addr,
	)

	return a.httpServer.ListenAndServe()

}

func (a *App) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down application")

	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("shutdown http server: %w", err)
	}

	if err := a.grpcClient.Close(); err != nil {
		return fmt.Errorf("close grpc clients: %w", err)
	}

	a.logger.Info("application stopped")

	return nil
}
