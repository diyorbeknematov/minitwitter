package handler

import (
	"log/slog"
	"net/http"

	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/service"
	"github.com/gorilla/websocket"
)

type Handler struct {
	service           *service.Service
	logger            *slog.Logger
	upgrader          websocket.Upgrader
	notificationWSURL string
}

func New(svc *service.Service, logger *slog.Logger, url string) *Handler {
	var upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	
	return &Handler{
		service:           svc,
		logger:            logger,
		upgrader:          upgrader,
		notificationWSURL: url,
	}
}
