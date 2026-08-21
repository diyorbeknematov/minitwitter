package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func (h *Handler) NotificationWS(c *gin.Context) {
	ctx := c.Request.Context()

	userID, ok := getUserID(c)
	if !ok {
		h.logger.WarnContext(ctx, "unauthorized websocket connection")

		errorResponse(c, http.StatusUnauthorized, "unauthorized")
		return
	}

	clientConn, err := h.upgrader.Upgrade(
		c.Writer,
		c.Request,
		nil,
	)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to upgrade websocket connection",
			"error", err,
		)
		return
	}

	defer clientConn.Close()

	header := http.Header{}
	header.Set("X-User-ID", userID.String())

	notificationConn, _, err := websocket.DefaultDialer.Dial(
		h.notificationWSURL,
		header,
	)
	if err != nil {
		h.logger.ErrorContext(ctx,
			"failed to connect to notification service",
			"error", err,
		)

		return
	}
	defer notificationConn.Close()

	// Client -> Notification Service
	go func() {
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				return
			}

			if err := notificationConn.WriteMessage(
				messageType,
				message,
			); err != nil {
				return
			}
		}
	}()

	// Notification Service → Client
	for {
		messageType, message, err :=
			notificationConn.ReadMessage()

		if err != nil {
			return
		}

		if err := clientConn.WriteMessage(
			messageType,
			message,
		); err != nil {
			return
		}
	}
}
