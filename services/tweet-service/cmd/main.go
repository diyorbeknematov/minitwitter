package cmd

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/cmd/app"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/logger"
)

func main() {
	cfg := config.Load()
	log := logger.SetupLog()

	a, err := app.New(cfg, log)
	if err != nil {
		log.Error("failed to create app", "error", err)
		os.Exit(1)
	}

	go func() {
		if err := a.Run(); err != nil {
			log.Error("gRPC server stopped", "error", err)
		}
	}()

	// Ctrl + C yoki docker stop kutamiz
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Info("shutting down...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.Shutdown(ctx); err != nil {
		log.Error("shotdown failed", "error", err)
	}
}
