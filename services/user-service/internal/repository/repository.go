package repository

import (
	"context"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository/postgres"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	UserRepo
	FollowRepo
	RefreshTokenRepo
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserRepo: postgres.NewUserRepo(db),
		FollowRepo: postgres.NewFollowRepo(db),
		RefreshTokenRepo: postgres.NewRefreshTokenRepo(db),
	}
}

type UserRepo interface {
	Create(context.Context, *models.User) error
	GetByID(context.Context, uuid.UUID) (*models.User, error)
	GetByEmail(context.Context, string) (*models.User, error)
	GetByUsername(context.Context, string) (*models.User, error)
	Search(context.Context, string, int, int) ([]models.User, int, error)
	GetUserFollowers(context.Context, uuid.UUID, int, int) ([]models.User, int, error)
	Update(context.Context, *models.User) error
	Delete(context.Context, uuid.UUID) error
}

type FollowRepo interface {
	Create(context.Context, *models.Follow) error
	Delete(context.Context, *models.Follow) error
	Exists(context.Context, uuid.UUID, uuid.UUID) (bool, error)
	GetFollowers(context.Context, models.GetFollowersReq) ([]models.Follow, int, error)
	GetFollowing(context.Context, models.GetFollowingReq) ([]models.Follow, int, error)
	CountFollowers(context.Context, uuid.UUID) (int, error)
	CountFollowing(context.Context, uuid.UUID) (int, error)
	GetFollowingIDs(context.Context, uuid.UUID) ([]uuid.UUID, error)
}

type RefreshTokenRepo interface {
	Create(context.Context, models.RefreshToken) error
	GetByUserID(context.Context, uuid.UUID) (models.RefreshToken, error)
	DeleteByUserID(context.Context, uuid.UUID) error
}
