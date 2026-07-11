package models

import (
	"time"

	"github.com/google/uuid"
)

type Follow struct {
	FollowerID  uuid.UUID `db:"follower_id"`
	FollowingID uuid.UUID `db:"following_id"`
	CreatedAt   time.Time `db:"created_at"`
}
