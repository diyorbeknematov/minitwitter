package dto

import (
	"time"

	"github.com/google/uuid"
)

type MediaCategory string

const (
	MediaCategoryTweet  MediaCategory = "tweet"
	MediaCategoryAvatar MediaCategory = "avatar"
)

type Media struct {
	ID              uuid.UUID `json:"id"`
	ObjectKey       string    `json:"object_key"`
	Url             string    `json:"url"`
	OriginalName    string    `json:"original_name"`
	MimeType        string    `json:"mime_type"`
	Size            int64     `json:"size"`
	StorageProvider string    `json:"storage_provider"`
	CreatedAt       time.Time `json:"created_at"`
}

type UploadMediaReq struct {
	Category MediaCategory `form:"category"`
}

type GetMediasReq struct {
	MediaIDs []uuid.UUID `json:"media_ids"`
}

type GetMediasResp struct {
	Medias []Media `json:"medias"`
}
