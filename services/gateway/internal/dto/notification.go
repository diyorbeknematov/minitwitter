package dto

import (
	"time"

	"github.com/google/uuid"
)

type NotificationType string

const (
	NotificationTypeFollow  NotificationType = "follow"
	NotificationTypeLike    NotificationType = "like"
	NotificationTypeRetweet NotificationType = "retweet"
	NotificationTypeReply   NotificationType = "reply"
)

type Notification struct {
	ID        uuid.UUID        `json:"id"`
	ActorID   uuid.UUID        `json:"actor_id"`
	Type      NotificationType `json:"type"`
	EntityID  uuid.UUID        `json:"entity_id"`
	IsRead    bool             `json:"is_read"`
	CreatedAt time.Time        `json:"created_at"`
}

type GetNotificationsReq struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

type GetNotificationsResp struct {
	Notifications []Notification `json:"notifications"`
	Pagination    Pagination     `json:"pagination"`
}

type GetUnreadCountResp struct {
	Count  int64     `json:"count"`
}
