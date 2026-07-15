package postgres

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type tweetRepo struct {
	db *sqlx.DB
}

func NewTweetRepo(db *sqlx.DB) *tweetRepo {
	return &tweetRepo{
		db: db,
	}
}

func (r *tweetRepo) Create(ctx context.Context, tweet *models.Tweet) error {
	query := `
		INSERT INTO tweets (
			author_id,
			content,
			reply_to_tweet_id,
			updated_at
		) 
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at, updated_at;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		tweet.ID,
		tweet.Content,
		tweet.ReplyToTweetID,
	).Scan(
		&tweet.ID,
		&tweet.Content,
		&tweet.ReplyToTweetID,
	)
	if err != nil {
		return apperror.Wrap("repository", "CreateTweet", "failed to create tweet", err)
	}

	return nil
}

func (r *tweetRepo) GetByID(ctx context.Context, tweetID string) (*models.Tweet, error) {
	query := `
		SELECT
			id,
			author_id,
			reply_to_tweet_id,
			created_at,
			updated_at
		FROM tweets
		WHERE id = $1;
	`

	var tweet models.Tweet
	err := r.db.QueryRow(query, tweetID).Scan(
		&tweet.ID,
		&tweet.AuthorID,
		&tweet.ReplyToTweetID,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
	)
	if err != nil {
		return nil, apperror.Wrap("repository", "GetTweetByID", "failed to get tweet by id", err)
	}

	return &tweet, nil
}
func (r *tweetRepo) GetByUser(ctx context.Context, userID string) (*models.Tweet, error)
func (r *tweetRepo) Retweet(ctx context.Context, tweetID string) error
func (r *tweetRepo) UndoRetweet(ctx context.Context, tweetID string) error
func (r *tweetRepo) Update(ctx context.Context) error
func (r *tweetRepo) Delete(ctx context.Context, tweetID string) error
func (r *tweetRepo) GetTimeline(ctx context.Context, getTimeline models.GetTimelineReq) (*models.Tweet, int, error)
