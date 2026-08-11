package dto

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID             uuid.UUID `json:"id"`
	Username       string    `json:"username"`
	Email          string    `json:"email"`
	Name           string    `json:"name"`
	Bio            string    `json:"bio"`
	AvatarURL      string    `json:"avatar_url"`
	FollowersCount uint64    `json:"followers_count"`
	FollowingCount uint64    `json:"following_count"`
	CreatedAt      time.Time `json:"created_at"`
}

type UsersResp struct {
	Users      []User     `json:"users"`
	Pagination Pagination `json:"pagination"`
}

type UpdateProfileReq struct {
	Name      string `json:"name"`
	Bio       string `json:"bio"`
	AvatarURL string `json:"avatar_url"`
}

type SearchUsersQuery struct {
	Query string `form:"q"`
	Page  int32  `form:"page"`
	Limit int32  `form:"limit"`
}
