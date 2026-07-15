package models

import (
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID             uuid.UUID  `json:"id" db:"id"`
	AuthorID       uuid.UUID  `json:"author_id" db:"author_id"`
	Content        string     `json:"content" db:"content"`
	ReplyToTweetID *uuid.UUID `json:"reply_to_tweet_id" db:"reply_to_tweet_id"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
	DeletedAt      *time.Time `db:"deleted_at"`
}

type Like struct {
	TweetID   uuid.UUID `json:"tweet_id" db:"tweet_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type Retweet struct {
	TweetID   uuid.UUID `json:"tweet_id" db:"tweet_id"`
	UserID    uuid.UUID `json:"user_id" db:"user_id"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}

type GetTweetsByUserReq struct {
	UserID uuid.UUID `db:"user_id"`
	Limit  int       `db:"limit"`
	Offset int       `db:"offset"`
}

type GetTimelineReq struct {
	UserIDs []uuid.UUID `db:"user_ids"`
	Limit   int         `db:"limit"`
	Offset  int         `db:"offset"`
}
