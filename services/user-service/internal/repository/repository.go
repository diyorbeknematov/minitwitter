package repository

import (
	"context"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	User
	Follow
	RefreshToken
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		User:         postgres.NewUserRepo(db),
		Follow:       postgres.NewFollowRepo(db),
		RefreshToken: postgres.NewRefreshTokenRepo(db),
	}
}

type User interface {
	Create(context.Context, *models.User) error
	GetByID(context.Context, uuid.UUID) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetByUsername(context.Context, string) (*models.User, error)
	Search(context.Context, string, int, int) ([]models.User, int, error)
	GetUserFollowers(context.Context, uuid.UUID, int, int) ([]models.User, int, error)
	GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.User, int, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, uuid.UUID) error
}

type Follow interface {
	Create(context.Context, *models.Follow) error
	Delete(context.Context, *models.Follow) error
	Exists(context.Context, uuid.UUID, uuid.UUID) (bool, error)
	GetFollowers(context.Context, models.GetFollowersReq) ([]models.Follow, int, error)
	GetFollowing(context.Context, models.GetFollowingReq) ([]models.Follow, int, error)
	CountFollowers(context.Context, uuid.UUID) (int, error)
	CountFollowing(context.Context, uuid.UUID) (int, error)
	GetFollowingIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
}

type RefreshToken interface {
	Create(context.Context, models.RefreshToken) error
	GetByUserID(context.Context, uuid.UUID) (models.RefreshToken, error)
	DeleteByUserID(context.Context, uuid.UUID) error
}
