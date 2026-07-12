package postgres

import (
	"context"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type refreshTokenRepo struct {
	db *sqlx.DB
}

func NewRefreshTokenRepo(db *sqlx.DB) *refreshTokenRepo {
	return &refreshTokenRepo{
		db: db,
	}
}

func (r *refreshTokenRepo) Create(ctx context.Context, token models.RefreshToken) error {
	query := `
		INSERT INTO refresh_token (
			id,
			user_id,
			token_hash,
			expires_at
		)
		VALUES ($1, $2, $3, $4);
	`

	_, err := r.db.Exec(query,
		token.ID,
		token.UserID,
		token.TokenHash,
		token.ExpiresAt,
	)

	return apperror.Wrap("repository", "RefreshTokenCreate", "failed to create refresh token", err)
}

func (r *refreshTokenRepo) GetByUserID(ctx context.Context, userID uuid.UUID) (models.RefreshToken, error) {
	query := `
	SELECT 
		id,
		user_id,
		token_hash,
		expires_at,
		created_at
	FROM refresh_tokens
	WHERE user_id = $1;
	`

	var refreshToken models.RefreshToken
	err := r.db.QueryRow(query, userID).Scan(
		&refreshToken.ID,
		&refreshToken.UserID,
		&refreshToken.TokenHash,
		&refreshToken.CreatedAt,
	)

	if err != nil {
		return models.RefreshToken{}, apperror.Wrap("repository", "GetTokenByUserID", "failed to get refresh token by user id", err)
	}

	return refreshToken, nil
}

func (r *refreshTokenRepo) DeleteByUserID(ctx context.Context, userID uuid.UUID) error {
	query := `
		DELETE FROM refresh_tokens
		WHERE user_id = $1;
	`

	_, err := r.db.Exec(query, userID)

	if err != nil {
		return apperror.Wrap("repository", "DeleteTokenBuUserID", "fialed to delete refresh token", err)
	}

	return nil
}
