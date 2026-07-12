package postgres

import (
	"context"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
)

type followRepo struct {
	db *sqlx.DB
}

func NewFollowRepo(db *sqlx.DB) *followRepo {
	return &followRepo{
		db: db,
	}
}

func (r *followRepo) Create(ctx context.Context, follow *models.Follow) error {
	query := `
		INSERT INTO follows (
			follower_id,
			following_id
		) 
		VALUES ($1, $2)
		RETURNING created_at;
	`
	err := r.db.QueryRowContext(
		ctx,
		query,
		follow.FollowerID,
		follow.FollowingID,
		follow.CreatedAt,
	).Scan(&follow.CreatedAt)

	if err != nil {
		return apperror.Wrap("repository", "UserFollowing", "failed to following", err)
	}

	return nil
}

func (r *followRepo) Delete(ctx context.Context, follow *models.Follow) error {
	query := `
		DELETE FROM fallows
		WHERE follower_id = $1 AND following_id = $2;
	`

	_, err := r.db.Exec(query,
		follow.FollowerID,
		follow.FollowingID,
	)

	if err != nil {
		return apperror.Wrap("repository", "UserUnfollowing", "failed to unfollowing", err)
	}

	return nil
}

func (r *followRepo) Exists(ctx context.Context, followerID, followingID uuid.UUID) (bool, error) {
	query := `
		SELECT EXISTS (
			SELECT 
				1 
			FROM follows
			WHERE follower_id = $1 AND following_id = $2
		);
	`
	var exist bool

	err := r.db.QueryRow(query, followerID, followingID).Scan(&exist)

	if err != nil {
		return exist, apperror.Wrap("repository", "CheckFollowing", "failed to check following", err)
	}

	return exist, nil
}

func (r *followRepo) GetFollowers(ctx context.Context, req models.GetFollowersReq) ([]models.Follow, int, error) {
	baseQuery := `
		SELECT 
			follower_id,
			following_id,
			created_at
		FROM follows
		WHERE following_id = $1 
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3;
	`
	countQuery := `SELECT COUNT(*) FROM follows WHERE following_id = $1`

	follows := []models.Follow{}
	rows, err := r.db.Query(baseQuery,
		req.UserID,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowers", "failed to execute query", err)
	}

	for rows.Next() {
		follow := models.Follow{}
		if err = rows.Scan(
			&follow.FollowerID,
			&follow.FollowingID,
			&follow.CreatedAt,
		); err != nil {
			return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowers", "failed to scan rows", err)
		}

		follows = append(follows, follow)
	}

	if err = rows.Err(); err != nil {
		return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowers", "failed to get rows", err)
	}

	var total int
	err = r.db.QueryRow(countQuery, req.UserID).Scan(&total)
	if err != nil {
		return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowers", "failed to get total count", err)
	}

	return follows, total, nil
}

func (r *followRepo) GetFollowing(ctx context.Context, req models.GetFollowingReq) ([]models.Follow, int, error) {
	baseQuery := `
		SELECT 
			follower_id,
			following_id,
			created_at
		FROM follows
		WHERE follower_id = $1
		ORDER BY created_at DESC 
		LIMIT $2 OFFSET $3;
	`
	countQuery := `SELECT COUNT(*) FROM follows WHERE follower_id = $1`

	follows := []models.Follow{}
	rows, err := r.db.Query(baseQuery,
		req.UserID,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowing", "failed to execute query", err)
	}

	for rows.Next() {
		follow := models.Follow{}
		if err = rows.Scan(
			&follow.FollowerID,
			&follow.FollowingID,
			&follow.CreatedAt,
		); err != nil {
			return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowing", "failed to scan rows", err)
		}

		follows = append(follows, follow)
	}

	if err = rows.Err(); err != nil {
		return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowing", "failed to get rows", err)
	}

	var total int
	err = r.db.QueryRow(countQuery, req.UserID).Scan(&total)
	if err != nil {
		return []models.Follow{}, 0, apperror.Wrap("repository", "GetFollowing", "failed to get total count", err)
	}

	return follows, total, nil
}

func (r *followRepo) CountFollowers(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT 
			COUNT(*)
		FROM follows
		WHERE following_id = $1;
	`

	var countFollowers int
	err := r.db.QueryRow(query, userID).Scan(&countFollowers)

	return countFollowers, apperror.Wrap("repository", "GetCountFollowers", "failed to get followers count", err)
}

func (r *followRepo) CountFollowing(ctx context.Context, userID uuid.UUID) (int, error) {
	query := `
		SELECT 
			COUNT(*)
		FROM follows
		WHERE follower_id = $1;
	`

	var countFollowing int
	err := r.db.QueryRow(query, userID).Scan(&countFollowing)

	return countFollowing, apperror.Wrap("repository", "GetCountFollowing", "failed to get following count", err)
}

func (r *followRepo) GetFollowingIDs(ctx context.Context, userID uuid.UUID) ([]uuid.UUID, error) {
	query := `
		SELECT 
			following_id
		FROM follows
		WHERE follower_id = $1; 
	`
	var ids []uuid.UUID

	rows, err := r.db.Query(query, userID)
	if err != nil {
		return []uuid.UUID{}, apperror.Wrap("repository", "GetFollowingIDs", "failed to get following ids", err)
	}
	defer rows.Close()

	for rows.Next() {
		var id uuid.UUID
		if err = rows.Scan(&id); err != nil {
			return []uuid.UUID{}, apperror.Wrap("repository", "GetFollowingIDs", "failed to scan rows", err)
		}

		ids = append(ids, id)
	}

	if err = rows.Err(); err != nil {
		return []uuid.UUID{}, apperror.Wrap("repository", "GetFollowingIDs", "failed to get rows", err)
	}

	return ids, nil
}
