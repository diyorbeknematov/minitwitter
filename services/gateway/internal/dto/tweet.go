package dto

import (
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID             uuid.UUID   `json:"id"`
	AuthorID       uuid.UUID   `json:"author_id"`
	Content        string      `json:"content"`
	MediaIDs       []uuid.UUID `json:"media_ids"`
	ReplyToTweetID *uuid.UUID   `json:"reply_to_tweet_id,omitempty"`
	LikesCount     int64       `json:"likes_count"`
	RetweetsCount  int64       `json:"retweets_count"`
	CreatedAt      time.Time   `json:"created_at"`
	UpdatedAt      time.Time   `json:"updated_at"`
}

type CreateTweetReq struct {
	Content        string      `json:"content"`
	MediaIDs       []uuid.UUID `json:"media_ids,"`
	ReplyToTweetID *uuid.UUID   `json:"reply_to_tweet_id,omitempty"`
}

type UpdateTweetReq struct {
	Content  string      `json:"content"`
	MediaIDs []uuid.UUID `json:"media_ids"`
}

type GetTweetByUserReq struct {
	Page  int32 `form:"page"`
	Limit int32 `form:"limit"`
}

type GetTweetByUserResp struct {
	Tweets     []Tweet    `json:"tweets"`
	Pagination Pagination `json:"pagination"`
}

type GetTimelineReq struct {
	UserIDs []uuid.UUID `json:"user_ids"`
	Page    int32       `form:"page"`
	Limit   int32       `form:"limit"`
}

type GetTimelineResp struct {
	Tweets     []Tweet    `json:"tweets"`
	Pagination Pagination `json:"pagination"`
}
