package ws

import (
	"errors"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO
	},
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}

	client := NewClient(userID, conn, h.hub)
	h.hub.Register(client)

	go client.ReadPump()
	go client.WritePump()
}

func getUserID(r *http.Request) (uuid.UUID, error) {
	userID := r.Header.Get("X-User-ID")
	if userID == "" {
		return uuid.Nil, errors.New("missing user id")
	}

	return uuid.Parse(userID)
}
