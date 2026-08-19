package handler

import (
	"log/slog"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/service"
)

type Handler struct {
	service *service.Service
	logger  *slog.Logger
}

func NewHandler(svc *service.Service, logger *slog.Logger) *Handler {
	return &Handler{
		service: svc,
		logger:  logger,
	}
}
