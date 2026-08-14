package models

import (
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID             uuid.UUID   `db:"id"`
	AuthorID       uuid.UUID   `db:"author_id"`
	Content        string      `db:"content"`
	ReplyToTweetID *uuid.UUID  `db:"reply_to_tweet_id"`

	MediaIDs       []uuid.UUID `db:"media_ids"`

	LikesCount     int64       `db:"likes_count"`
	RetweetsCount  int64       `db:"retweets_count"`
	
	CreatedAt      time.Time   `db:"created_at"`
	UpdatedAt      time.Time   `db:"updated_at"`
	DeletedAt      *time.Time  `db:"deleted_at"`
}

type Like struct {
	TweetID   uuid.UUID `db:"tweet_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}

type Retweet struct {
	TweetID   uuid.UUID `db:"tweet_id"`
	UserID    uuid.UUID `db:"user_id"`
	CreatedAt time.Time `db:"created_at"`
}
