package repository

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Notif
}

func NewRepo(db *sqlx.DB) *Repository {
	return &Repository{
		Notif: postgres.NewNotificationRepo(db),
	}
}

type Notif interface {
	Create(context.Context, *models.Notification) error
	GetByID(context.Context, string) (*models.Notification, error)
	GetMany(context.Context, models.GetNotifications) ([]models.Notification, int, error)
	Update(context.Context, string) error
	UpdateAll(context.Context, string) error
	Delete(context.Context, string) error
	GetUnreadCount(context.Context, string) (int, error)
}
