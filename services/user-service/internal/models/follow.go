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

type GetFollowersReq struct {
	UserID uuid.UUID `db:"user_id"`
	Limit  int       `db:"limit"`
	Offset int       `db:"offset"`
}

type GetFollowingReq struct {
	UserID uuid.UUID `db:"user_id"`
	Limit  int       `db:"limit"`
	Offset int       `db:"offest"`
}
