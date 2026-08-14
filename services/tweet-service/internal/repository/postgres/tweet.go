package postgres

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/google/uuid"
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
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap("repository", "CreateTweet", "failed to begin transaction", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO tweets (
			author_id,
			content,
			reply_to_tweet_id
		) 
		VALUES ($1, $2, $3)
		RETURNING id, created_at, updated_at;
	`
	err = tx.QueryRowContext(
		ctx,
		query,
		tweet.AuthorID,
		tweet.Content,
		tweet.ReplyToTweetID,
	).Scan(
		&tweet.ID,
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
	)
	if err != nil {
		return apperror.Wrap("repository", "CreateTweet", "failed to create tweet", err)
	}

	query = `
		INSERT INTO tweet_media (
			tweet_id,
			media_id,
			position
		)
		VALUES($1, $2, $3);
	`
	for i, mediaID := range tweet.MediaIDs {
		_, err := tx.ExecContext(ctx, query, tweet.ID, mediaID, i)
		if err != nil {
			return apperror.Wrap("repository", "CreateTweet", "failed to create tweet_media", err)
		}
	}

	if err := tx.Commit(); err != nil {
		return apperror.Wrap(
			"repository", "CreateTweet", "failed to commit transaction", err,
		)
	}

	return nil
}

func (r *tweetRepo) GetByID(ctx context.Context, tweetID string) (*models.Tweet, error) {
	query := `
		SELECT 
			t.id, 
			t.author_id, 
			t.content, 
			t.reply_to_tweet_id,
			t.likes_count,
			t.retweets_count,
			COALESCE(
				ARRAY_AGG(tm.media_id ORDER BY tm.position) FILTER (WHERE tm.media_id IS NOT NULL), 
				'{}'
				) AS media_ids,
			t.created_at,
			t.updated_at
		FROM tweets t
		LEFT JOIN tweet_media tm ON tm.tweet_id = t.id
		WHERE t.id = $1 AND t.deleted_at IS NULL
		GROUP BY t.id
		ORDER BY t.created_at DESC;
	`

	var tweet models.Tweet
	err := r.db.QueryRow(query, tweetID).Scan(
		&tweet.ID,
		&tweet.AuthorID,
		&tweet.Content,
		&tweet.ReplyToTweetID,
		&tweet.LikesCount,
		&tweet.RetweetsCount,
		pq.Array(&tweet.MediaIDs),
		&tweet.CreatedAt,
		&tweet.UpdatedAt,
	)
	if err != nil {
		return nil, apperror.Wrap("repository", "GetTweetByID", "failed to get tweet by id", err)
	}

	return &tweet, nil
}

func (r *tweetRepo) GetByUser(ctx context.Context, userID string, limit, offset int32) ([]models.Tweet, int, error) {
	baseQuery := `
		SELECT 
			id,
			author_id,
			content,
			reply_to_tweet_id,
			likes_count,
			retweets_count,
			COALESCE(
				ARRAY_AGG(tm.media_id ORDER BY tm.position) FILTER (WHERE tm.media_id IS NOT NULL), 
				'{}'
				) AS media_ids,
			created_at,
			updated_at
		FROM tweets t
		LEFT JOIN tweet_media tm ON tm.tweet_id = t.id
		WHERE author_id = $1 AND t.deleted_at IS NULL
		GROUP BY t.id 
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3;
	`
	countQuery := `
		SELECT 
			COUNT(*)
		FROM tweets
		WHERE author_id = $1 AND t.deleted_at IS NULL;
	`

	rows, err := r.db.QueryContext(ctx, baseQuery, userID, limit, offset)
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
			&tweet.LikesCount,
			&tweet.RetweetsCount,
			pq.Array(&tweet.MediaIDs),
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
	err = r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTweetsByUser", "failed to get total number", err)
	}

	return tweets, total, nil
}

func (r *tweetRepo) CreateRetweet(ctx context.Context, retweet models.Retweet) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap(
			"repository", "CreateRetweet", "failed to begin transaction", err,
		)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		INSERT INTO retweets (
			tweet_id,
			user_id
		)
		VALUES($1, $2)
	`

	_, err = tx.ExecContext(ctx,
		query,
		retweet.TweetID,
		retweet.UserID,
	)
	if err != nil {
		return apperror.Wrap("repository", "CreateRetweet", "failed to tweet to retweet", err)
	}

	query = `
		UPDATE tweets
		SET 
			retweets_count = retweets_count + 1
		WHERE id = $1
	`
	_, err = tx.ExecContext(ctx, query, retweet.TweetID)
	if err != nil {
		return apperror.Wrap("repository", "CreataeRetweet", "failed to update retweets count", err)
	}

	if err := tx.Commit(); err != nil {
		return apperror.Wrap(
			"repository", "CreateRetweet", "failed to commit transaction", err,
		)
	}

	return nil
}

func (r *tweetRepo) DeleteRetweet(ctx context.Context, tweetID, userID string) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap("repository", "DeleteRetweet", "failed to begin transaction", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()

	query := `
		DELETE FROM retweets
		WHERE tweet_id = $1 
			AND user_id = $2;
	`

	res, err := tx.ExecContext(ctx, query, tweetID, userID)
	if err != nil {
		return apperror.Wrap("repository", "DeleteRetweet", "failed to delete retweet", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "DeleteRetweet", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "DeleteRetweet", "no rows affected to delete retweet", err)
	}

	query = `
		UPDATE tweets
		SET 
			retweets_count = retweets_count - 1
		WHERE id = $1
	`

	_, err = tx.ExecContext(ctx, tweetID)
	if err != nil {
		apperror.Wrap("repository", "DeleteRetweet", "failed to update retweets count", err)
	}

	if err := tx.Commit(); err != nil {
		return apperror.Wrap("repository", "DeleteRetweet", "failed to commit transaction", err)
	}

	return nil
}

func (r *tweetRepo) Update(ctx context.Context, tweet models.Tweet) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return apperror.Wrap("repository", "UpdateTweet", "failed to begin transaction", err)
	}
	defer tx.Rollback()

	query := `
		UPDATE tweets 
		SET
			content = $2,
			updated_at = $3
		WHERE tweet_id = $1
	`
	res, err := tx.Exec(query,
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

func (r *tweetRepo) GetTimeline(ctx context.Context, userIDs []uuid.UUID, limit, offset int32) ([]models.Tweet, int, error) {
	baseQuery := `
		SELECT 
			t.id, 
			t.author_id, 
			t.content, 
			t.reply_to_tweet_id,
			t.likes_count,
			t.retweets_count,
			COALESCE(
				ARRAY_AGG(tm.media_id ORDER BY tm.position) FILTER (WHERE tm.media_id IS NOT NULL), 
				'{}'
				) AS media_ids,
			t.created_at,
			t.updated_at
		FROM tweets t
		LEFT JOIN tweet_media tm ON tm.tweet_id = t.id
		WHERE t.author_id = ANY($1) AND t.deleted_at IS NULL
		GROUP BY t.id
		ORDER BY t.created_at DESC
		LIMIT $2 OFFSET $3;
	`
	countQuery := `
		SELECT
			COUNT(*)
		FROM tweets
		WHERE deleted_at IS NULL AND author_id = ANY($1);
	`

	rows, err := r.db.QueryContext(
		ctx,
		baseQuery,
		pq.Array(userIDs),
		limit,
		offset,
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
			&tweet.LikesCount,
			&tweet.RetweetsCount,
			pq.Array(&tweet.MediaIDs),
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
	err = r.db.QueryRowContext(ctx, countQuery, pq.Array(userIDs)).Scan(&total)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetTimeline", "failed to get total count", err)
	}

	return tweets, total, nil
}
