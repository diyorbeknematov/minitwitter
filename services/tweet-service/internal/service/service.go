package service

import (
	"log/slog"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/repository"
)

type Service struct {
	Tweet *tweetService
}

func NewService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		Tweet: NewTweetService(repo, logger),
	}
}
