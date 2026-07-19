package postgres

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type likeRepo struct {
	db *sqlx.DB
}

func NewLikeRepo(db *sqlx.DB) *likeRepo {
	return &likeRepo{
		db: db,
	}
}

func (r *likeRepo) Create(ctx context.Context, like *models.Like) error {
	query := `
		INSERT INTO likes (
			tweet_id,
			user_id,
			created_at = $3
		) 
		VALUES($1, $2)
	`

	_, err := r.db.ExecContext(ctx, query, like.TweetID, like.UserID, like.CreatedAt)
	if err != nil {
		return apperror.Wrap("repository", "CreateLike", "failed to create like", err)
	}

	return nil
}

func (r *likeRepo) Delete(ctx context.Context, tweetID, userID string) error {
	query := `
		DELETE FROM likes
		WHERE tweet_id = $1
			AND user_id = $2;
	`

	res, err := r.db.ExecContext(ctx,
		query,
		tweetID,
		userID,
	)
	if err != nil {
		return apperror.Wrap("repository", "Unlike", "failed to unlike", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "Unlike", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "Unlike", "no rows affected to unlike", err)
	}

	return nil
}
