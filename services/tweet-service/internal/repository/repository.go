package repository

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/repository/postgres"
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
	Create(ctx context.Context, tweet *models.Tweet) error
	GetByID(ctx context.Context, tweetID string) (*models.Tweet, error)
	GetByUser(ctx context.Context, userID string, limit, offset int) ([]models.Tweet, int, error)
	GetTimeline(ctx context.Context, req models.GetTimelineReq) ([]models.Tweet, int, error)
	CreateRetweet(ctx context.Context, retweet models.Retweet) error
	DeleteRetweet(ctx context.Context, tweetID, userID string) error
	Update(ctx context.Context, tweet models.Tweet) error
	Delete(ctx context.Context, tweetID string) error
}

type Like interface {
	Create(ctx context.Context, like *models.Like) error
	Delete(ctx context.Context, tweetID, userID string) error
}
