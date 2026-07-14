package service

import (
	"log/slog"

	"github.com/diyorbek/minitwitter/services/user-service/internal/config"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository"
)

type Service struct {
	Auth *authService
	User *userService
}

func NewService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *Service {
	return &Service{
		Auth: NewAuthService(repo, cfg, logger),
		User: NewUserService(repo, logger),
	}
}
