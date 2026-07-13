package service

import (
	"log/slog"

	"github.com/diyorbek/minitwitter/services/user-service/internal/repository"
)

type userSerivice struct {
	repo   *repository.Repository
	logger *slog.Logger
}

func NewUserService(repo *repository.Repository, logger *slog.Logger) *userSerivice {
	return &userSerivice{
		repo:   repo,
		logger: logger,
	}
}

