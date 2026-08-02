package service

import (
	"log/slog"

	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/storage"
)

type Service struct {
	Media *mediaService
}

func NewService(repo *repository.Repository, store storage.ObjectStorage, logger *slog.Logger, cfg *config.Config) *Service {
	return &Service{
		Media: NewMediaService(repo, store, logger, cfg),
	}
}
