package models

import (
	"time"

	"github.com/google/uuid"
)

type Media struct {
	ID              uuid.UUID `json:"id" db:"id"`
	OwnerID         uuid.UUID `json:"owner_id" db:"owner_id"`
	ObjectKey       string    `json:"object_key" db:"object_key"`
	OriginalName    string    `json:"original_name" db:"original_name"`
	MimeType        string    `json:"mime_type" db:"mime_type"`
	Size            uint64    `json:"size" db:"size"`
	StorageProvider string    `json:"storage_provider" db:"storage_provider"`
	CreatedAt       time.Time `json:"created_at" db:"created_at"`
}
