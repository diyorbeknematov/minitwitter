package service

import (
	"log/slog"

	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/repository"
)

type Service struct {
	Notification *notificationService
}

func NewService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		Notification: NewNotificationService(repo, cfg, logger),
	}
}
