package repository

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/repository/postgres"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Media
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Media: postgres.NewMediaRepo(db),
	}
}

type Media interface {
	Create(context.Context, *models.Media) error
	GetByID(context.Context, string) (*models.Media, error)
	GetMany(context.Context, []string) ([]models.Media, error)
	Delete(context.Context, string) error
}
