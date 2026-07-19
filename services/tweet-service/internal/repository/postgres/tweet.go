package postgres

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
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
func (r *tweetRepo) GetByUser(ctx context.Context, userID string, limit, offset int) ([]models.Tweet, int, error) {
	baseQuery := `
		SELECT 
			id,
			author_id,
			content,
			reply_to_tweet_id,
			created_at,
			updated_at
		FROM tweets
		WHERE author_id = $1
		LIMIT $2 OFFSET $3;
	`
	countQuery := `
		SELECT 
			COUNT(*)
		FROM tweets
		WHERE author_id = $1;
	`

	rows, err := r.db.Query(baseQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTweetByUser", "failed to get tweet by user id", err)
	}
	defer rows.Close()

	tweets := []models.Tweet{}
	for rows.Next() {
		tweet := models.Tweet{}
		if err = rows.Scan(
			&tweet.ID,
			&tweet.AuthorID,
			&tweet.Content,
			&tweet.ReplyToTweetID,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
		); err != nil {
			return nil, 0, apperror.Wrap("repository", "GetTweetsByUser", "failed to scan tweet", err)
		}

		tweets = append(tweets, tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTweetsByUser", "failed to check rows error", err)
	}

	var total int
	err = r.db.QueryRow(countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTweetsByUser", "failed to get total number", err)
	}

	return tweets, total, nil
}

func (r *tweetRepo) Retweet(ctx context.Context, retweet models.Retweet) error {
	query := `
		INSERT INTO retweets (
			tweet_id,
			user_id,
			created_at
		)
		VALUES($1, $2, $3)
	`

	_, err := r.db.Exec(query,
		retweet.TweetID,
		retweet.UserID,
		retweet.CreatedAt,
	)
	if err != nil {
		return apperror.Wrap("repository", "Retweet", "failed to tweet to retweet", err)
	}

	return nil
}

func (r *tweetRepo) UndoRetweet(ctx context.Context, tweetID string) error {
	query := `
		DELETE FROM retweets
		WHERE tweet_id = $1
	`

	res, err := r.db.Exec(query, tweetID)
	if err != nil {
		return apperror.Wrap("repository", "UndoRetweet", "failed to undoretweet a tweet", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "UndoRetweet", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "Undoretweet", "no rows affected to undoretweet", err)
	}

	return nil
}

func (r *tweetRepo) Update(ctx context.Context, tweet models.Tweet) error {
	query := `
		UPDATE tweets 
		SET
			content = $2,
			updated_at = $3
		WHERE tweet_id = $1
	`
	res, err := r.db.Exec(query,
		tweet.ID,
		tweet.Content,
		tweet.UpdatedAt,
	)
	if err != nil {
		return apperror.Wrap("repository", "UpdateTweet", "failed to update tweet", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "UpdateTweet", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "UpdateTweet", "no rows affected to update tweet", err)
	}

	return nil
}

func (r *tweetRepo) Delete(ctx context.Context, tweetID string) error {
	query := `
		UPDATE tweets
		SET
			updated_at = now()
		WHERE id = $1
	`

	res, err := r.db.Exec(query, tweetID)
	if err != nil {
		return apperror.Wrap("repository", "DeleteTweet", "failed to delete tweet", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "DeleteTweet", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "DeleteTweet", "no rows affected to delete tweet", err)
	}

	return nil
}

func (r *tweetRepo) GetTimeline(ctx context.Context, req models.GetTimelineReq) ([]models.Tweet, int, error) {
	baseQuery := `
		SELECT
			id,
			author_id,
			content,
			reply_to_tweet_id,
			created_at,
			updated_at
		FROM tweets
		WHERE author_id = ANY($1)
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3;
	`
	countQuery := `
		SELECT
			COUNT(*)
		FROM tweets
		WHERE author_id = ANY($1);
	`

	rows, err := r.db.QueryContext(
		ctx,
		baseQuery,
		pq.Array(req.UserIDs),
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTimeline", "failed to execute the query context", err)
	}
	defer rows.Close()

	tweets := []models.Tweet{}
	for rows.Next() {
		tweet := models.Tweet{}
		if err = rows.Scan(
			&tweet.ID,
			&tweet.AuthorID,
			&tweet.Content,
			&tweet.ReplyToTweetID,
			&tweet.CreatedAt,
			&tweet.UpdatedAt,
		); err != nil {
			return nil, 0, apperror.Wrap("repository", "GetTimeline", "failed to scan tweet", err)
		}

		tweets = append(tweets, tweet)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperror.Wrap("respository", "GetTimeline", "failed to check rows to error", err)
	}

	var total int
	err = r.db.QueryRowContext(ctx, countQuery, pq.Array(req.UserIDs)).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTimeline", "failed to get total count", err)
	}

	return tweets, total, nil
}
