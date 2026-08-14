package repository

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	Tweet
	Like
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		Tweet: postgres.NewTweetRepo(db),
		Like:  postgres.NewLikeRepo(db),
	}
}

type Tweet interface {
	Create(context.Context, *models.Tweet) error
	GetByID(context.Context, string) (*models.Tweet, error)
	GetByUser(context.Context, string, int32, int32) ([]models.Tweet, int, error)
	GetTimeline(context.Context, []uuid.UUID, int32, int32) ([]models.Tweet, int, error)
	CreateRetweet(context.Context, models.Retweet) error
	DeleteRetweet(context.Context, string, string) error
	Update(context.Context, models.Tweet) error
	Delete(context.Context, string) error
}

type Like interface {
	Create(context.Context, *models.Like) error
	Delete(context.Context, string, string) error
}
