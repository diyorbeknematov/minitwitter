package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/diyorbeknematov/minitwitter/services/notification-service/cmd/app"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/pkg/logger"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()
	log := logger.SetupLog()

	a, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	if err := a.Run(ctx); err != nil {
		log.Error("app stopped with error", "errror", err)
		os.Exit(1)
	}

	log.Info("app stopped gracefully")

}
