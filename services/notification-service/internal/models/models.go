package models

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationFollow  NotificationType = "FOLLOW"
	NotificationLike    NotificationType = "LIKE"
	NotificationRetweet NotificationType = "RETWEET"
	NotificationReply   NotificationType = "REPLY"
)

type Notification struct {
	ID        uuid.UUID        `json:"id" db:"id"`
	UserID    uuid.UUID        `json:"user_id" db:"user_id"`
	ActorID   uuid.UUID        `json:"actor_id" db:"actor_id"`
	Type      NotificationType `json:"type" db:"type"`
	EntityID  uuid.UUID        `json:"entity_id" db:"entity_id"`
	IsRead    bool             `json:"is_read" db:"is_read"`
	CreatedAt time.Time        `json:"created_at" db:"created_at"`
}

type GetNotifications struct {
	UserID uuid.UUID `json:"user_id" db:"user_id"`
	Limit  int       `json:"limit" db:"limit"`
	Offset int       `json:"offset" db:"offset"`
}
