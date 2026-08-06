package ws

import (
	"encoding/json"
	"sync"

	"github.com/google/uuid"
)

type Hub struct {
	clients map[uuid.UUID]*Client
	mu      sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[uuid.UUID]*Client),
	}
}

func (h *Hub) Register(c *Client) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.clients[c.UserID] = c
}

func (h *Hub) Unregister(userID uuid.UUID) {
	h.mu.Lock()
	defer h.mu.Unlock()

	if client, ok := h.clients[userID]; ok {
		close(client.Send)
		client.Conn.Close()
		delete(h.clients, userID)
	}
}

func (h *Hub) Send(userID uuid.UUID, v any) {
	h.mu.RLock()

	client, ok := h.clients[userID]

	h.mu.RUnlock()

	if !ok {
		return
	}

	payload, err := json.Marshal(v)
	if err != nil {
		return
	}

	select {
	case client.Send <- payload:
	default:
		h.Unregister(client.UserID)
	}
}

func (h *Hub) Broadcast(v any) {
	payload, err := json.Marshal(v)
	if err != nil {
		return
	}

	h.mu.RLock()
	var stale []uuid.UUID

	for id, client := range h.clients {
		select {
		case client.Send <- payload:
		default:
			stale = append(stale, id)
		}
	}
	h.mu.RUnlock()

	for _, id := range stale {
		h.Unregister(id)
	}
}

func (h *Hub) Close() {
	h.mu.Lock()
	defer h.mu.Unlock()

	for _, client := range h.clients {
		client.Conn.Close()
	}

	h.clients = make(map[uuid.UUID]*Client)
}
