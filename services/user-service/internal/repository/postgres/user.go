package postgres

import (
	"context"
	"strings"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) *userRepo {
	return &userRepo{
		db: db,
	}
}

func (r *userRepo) Create(ctx context.Context, user *models.User) error {
	query := `
        INSERT INTO users (
            username,
            email,
            password_hash,
            name
        )
        VALUES ($1, $2, $3, $4)
        RETURNING id, created_at, updated_at;
    `

	err := r.db.QueryRowContext(
		ctx,
		query,
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Name,
	).Scan(
		&user.ID,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return apperror.Wrap("repository", "CreateUser", "failed to create user", err)
	}

	return nil
}

func (r *userRepo) GetByID(ctx context.Context, id uuid.UUID) (*models.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			name,
			bio,
			avatar_media_id
		FROM users
		WHERE id = $1 AND deleted_at IS NULL;
	`

	var user models.User

	err := r.db.QueryRow(query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Bio,
		user.AvatarMediaID,
	)
	if err != nil {
		return nil, apperror.Wrap("repository", "GetUserByID", "failed to get user by id", err)
	}

	return &user, nil
}

func (r *userRepo) GetByEmail(ctx context.Context, email string) (*models.User, error) {
	query := `
		SELECT 
			id,
			email,
			username,
			password_hash
		FROM users
		WHERE email = $1 AND deleted_at IS NULL;
	`

	var user models.User

	if err := r.db.QueryRow(query, email).Scan(
		&user.ID,
		&user.Email,
		&user.Username,
		&user.PasswordHash,
	); err != nil {
		return nil, apperror.Wrap("repository", "GetUserByEmail", "failed to get user by email", err)
	}

	return &user, nil
}

func (r *userRepo) GetByUsername(ctx context.Context, username string) (*models.User, error) {
	query := `
		SELECT
			id,
			username,
			email,
			name,
			bio,
			avatar_mdia_id,
			created_at
		FROM users
		WHERE deleted_at IS NULL AND username = $1;
	`

	var user models.User

	if err := r.db.QueryRow(query, username).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Name,
		&user.Bio,
		user.AvatarMediaID,
		&user.CreatedAt,
	); err != nil {
		return nil, apperror.Wrap("repository", "GetUserByUsername", "failed to get user by username", err)
	}

	return &user, nil
}

func (r *userRepo) Search(ctx context.Context, search string, offset, limit int) ([]models.User, int, error) {
	baseQuery := `
		SELECT 
			id,
			username,
			name,
			avatar_media_id
		FROM users
		WHERE deleted_at IS NULL 
	`
	countQuery := `SELECT COUNT(*) FROM users deleted_at IS NULL `
	conditions := []string{}

	params := map[string]any{
		"limit":  limit,
		"offset": offset,
	}

	// Add search condition
	if search != "" {
		conditions = append(conditions, "(name || username) ILIKE :search")
		params["search"] = "%" + search + "%"
	}

	// Add WHERE clause if conditons exist
	if len(conditions) > 0 {
		whereClause := " AND " + strings.Join(conditions, " AND ")
		baseQuery += whereClause
		countQuery += whereClause
	}

	// Add pagination
	baseQuery += " LIMIT :limit OFFSET :offset"

	users := []models.User{}
	rows, err := r.db.NamedQuery(baseQuery, params)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "SearchUser", "failed to execute named query", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Name,
			user.AvatarMediaID,
		); err != nil {
			return nil, 0, apperror.Wrap("repository", "SearchUser", "failed to scan row", err)
		}

		users = append(users, user)
	}

	// Execute the count query
	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "SearchUser", "failed to execute named query", err)
	}
	countQuery = r.db.Rebind(countQuery)

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, apperror.Wrap("repository", "SearchUser", "failed to get total count", err)
	}

	return users, total, nil
}

func (r *userRepo) GetUserFollowers(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.User, int, error) {
	baseQuery := `
		SELECT 
			u.id,
			u.username,
			u.name,
			u.avatar_media_id
		FROM users AS u
		INNER JOIN follows AS f 
			ON u.id = f.follower_id
		WHERE f.following_id = $1
		LIMIT $2 OFFSET $3;
	`
	countQuery := `
		SELECT
			COUNT(*)
		FROM users AS u
		INNER JOIN follows AS f
			ON u.id = f.follower_id
		WHERE f.following_id = $1
	`

	var users []models.User
	rows, err := r.db.Query(baseQuery, userID, limit, offset)
	if err != nil {
		return []models.User{}, 0, apperror.Wrap("repository", "GetUserFollers", "failed to get followers", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User
		if err = rows.Scan(
			&user.ID,
			&user.Username,
			&user.Name,
			user.AvatarMediaID,
		); err != nil {
			return []models.User{}, 0, apperror.Wrap("repository", "GetUserFallowers", "failed to scan rows", err)
		}

		users = append(users, user)
	}

	if err = rows.Err(); err != nil {
		return []models.User{}, 0, apperror.Wrap("repository", "GetUserFollowers", "failed to get users", err)
	}

	var total int
	if err = r.db.QueryRow(countQuery, userID).Scan(&total); err != nil {
		return []models.User{}, 0, apperror.Wrap("repository", "GetUserFollowers", "failed to get total count", err)
	}

	return users, total, nil
}

func (r *userRepo) GetUserFollowing(ctx context.Context, userID uuid.UUID, limit, offset int) ([]models.User, int, error) {
	baseQuery := `
		SELECT
			u.id,
			u.username,
			u.name,
			u.avatar_media_id
		FROM users AS u
		INNER JOIN follows AS f
			ON u.id = f.following_id
		WHERE f.follower_id = $1
		LIMIT $2 OFFSET $3;
	`

	countQuery := `
		SELECT
			COUNT(*)
		FROM users AS u
		INNER JOIN follows AS f
			ON u.id = f.following_id
		WHERE f.follower_id = $1;
	`

	var users []models.User

	rows, err := r.db.QueryContext(ctx, baseQuery, userID, limit, offset)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetUserFollowing", "failed to get following users", err)
	}
	defer rows.Close()

	for rows.Next() {
		var user models.User

		if err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Name,
			&user.AvatarMediaID,
		); err != nil {
			return nil, 0, apperror.Wrap("repository", "GetUserFollowing", "failed to scan row", err)
		}

		users = append(users, user)
	}

	if err := rows.Err(); err != nil {
		return nil, 0, apperror.Wrap("repository", "GetUserFollowing", "rows iteration failed", err)
	}

	var total int
	if err := r.db.QueryRowContext(ctx, countQuery, userID).Scan(&total); err != nil {
		return nil, 0, apperror.Wrap("repository", "GetUserFollowing", "failed to get following count", err)
	}

	return users, total, nil
}

func (r *userRepo) Update(ctx context.Context, user *models.User) error {
	query := `
		UPDATE users
		SET
			username = $2,
			email = $3,
			name = $4,
			bio = $5,
			avatar_media_id = $6,
			updated_at = $7
		WHERE deleted_at IS NULL AND id = $1;
	`
	res, err := r.db.Exec(query,
		user.ID,
		user.Username,
		user.Email,
		user.Name,
		user.Bio,
		user.AvatarMediaID,
		user.UpdatedAt,
	)
	if err != nil {
		return apperror.Wrap("repository", "UserUpdate", "failed to update user", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "UserUpdate", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "UserUpdate", "no rows affected to user update", err)
	}

	return nil
}

func (r *userRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		UPDATE users
		SET
			deleted_at = $2
		WHERE deleted_at IS NULL AND id = $1;
	`

	res, err := r.db.Exec(query, id)
	if err != nil {
		return apperror.Wrap("repository", "UserDelete", "failed to delete user", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "UserDelete", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "UserDelete", "no rows effected on user delete", err)
	}

	return nil
}
