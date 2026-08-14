package dto

import (
	"time"

	"github.com/google/uuid"
)

type Tweet struct {
	ID      uuid.UUID `json:"id"`
	Content string    `json:"content"`

	AuthorID uuid.UUID `json:"-"` // ichkarida ishlatamiz
	Author   User      `json:"author"`

	MediaIDs  []uuid.UUID `json:"-"` // ichkarida ishlatamiz
	MediaURLs []string    `json:"media_urls"`

	ReplyToTweetID *uuid.UUID `json:"reply_to_tweet_id,omitempty"`

	LikesCount    int64
	RetweetsCount int64

	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateTweetReq struct {
	Content        string      `json:"content"`
	MediaIDs       []uuid.UUID `json:"media_ids,"`
	ReplyToTweetID *uuid.UUID  `json:"reply_to_tweet_id,omitempty"`
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
