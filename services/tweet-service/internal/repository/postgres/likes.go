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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap(
			"repository", "CreateLike", "failed to begin transaction", err,
		)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO likes (
			tweet_id,
			user_id,
			created_at
		)
		VALUES ($1, $2, $3);
	`

	_, err = tx.ExecContext(
		ctx,
		query,
		like.TweetID,
		like.UserID,
		like.CreatedAt,
	)
	if err != nil {
		return apperror.Wrap(
			"repository", "CreateLike", "failed to create like", err,
		)
	}

	query = `
		UPDATE tweets
		SET likes_count = likes_count + 1
		WHERE id = $1;
	`

	_, err = tx.ExecContext(ctx, query, like.TweetID)
	if err != nil {
		return apperror.Wrap(
			"repository", "CreateLike", "failed to update likes count", err,
		)
	}

	if err := tx.Commit(); err != nil {
		return apperror.Wrap(
			"repository",
			"CreateLike",
			"failed to commit transaction",
			err,
		)
	}

	return nil
}

func (r *likeRepo) Delete(ctx context.Context, tweetID, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap(
			"repository", "DeleteLike", "failed to begin transaction", err,
		)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		DELETE FROM likes
		WHERE tweet_id = $1
		  AND user_id = $2;
	`

	res, err := tx.ExecContext(
		ctx,
		query,
		tweetID,
		userID,
	)
	if err != nil {
		return apperror.Wrap(
			"repository", "DeleteLike", "failed to delete like", err,
		)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap(
			"repository", "DeleteLike", "failed to get rows affected", err,
		)
	}

	if rows == 0 {
		return apperror.Wrap(
			"repository", "DeleteLike", "like not found", nil,
		)
	}

	query = `
		UPDATE tweets
		SET likes_count = likes_count - 1
		WHERE id = $1;
	`

	_, err = tx.ExecContext(ctx, query, tweetID)
	if err != nil {
		return apperror.Wrap(
			"repository", "DeleteLike", "failed to update likes count",
			err,
		)
	}

	if err := tx.Commit(); err != nil {
		return apperror.Wrap(
			"repository", "DeleteLike", "failed to commit transaction",
			err,
		)
	}

	return nil
}