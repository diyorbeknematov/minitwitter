package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diyorbeknematov/minitwitter/services/gateway/cmd/app"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/logger"
)

func main() {
	cfg := config.Load()

	log := logger.SetupLog()

	app, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := app.Run(); err != nil {
			log.Error("application stopped with error", "error", err)
		}
	}()

	// Ctrl+C yoki docker stop kutamiz
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := app.Shutdown(ctx); err != nil {
		log.Error("shutdown failed", "error", err)
	}
}
