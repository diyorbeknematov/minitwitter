package postgres

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type userRepo struct {
	db *sqlx.DB
}

func NewAuthRepo(db *sqlx.DB) *userRepo {
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

	return err
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
		return nil, err
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
		return nil, err
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
		return nil, err
	}

	return &user, nil
}

func (r *userRepo) Search(ctx context.Context, search string, page, limit int) ([]models.User, int, error) {
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
		"limit": limit,
		"page":  page,
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
		return nil, 0, err
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
			return nil, 0, err
		}

		users = append(users, user)
	}

	// Execute the count query
	var total int
	countQuery, countArgs, err := sqlx.Named(countQuery, params)
	if err != nil {
		return nil, 0, err
	}
	countQuery = r.db.Rebind(countQuery)

	if err := r.db.Get(&total, countQuery, countArgs...); err != nil {
		return nil, 0, err
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
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.New("failed to get rows affected: " + err.Error())
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected")
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
		return err
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return errors.New("failed to get rows affected: " + err.Error())
	}

	if rows == 0 {
		return fmt.Errorf("no rows affected")
	}

	return nil
}
