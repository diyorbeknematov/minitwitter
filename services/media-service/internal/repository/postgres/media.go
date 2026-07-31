package postgres

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/media-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/media-service/pkg/apperror"
	"github.com/jmoiron/sqlx"
	"github.com/lib/pq"
)

type mediaRepo struct {
	db *sqlx.DB
}

func NewMediaRepo(db *sqlx.DB) *mediaRepo {
	return &mediaRepo{
		db: db,
	}
}

func (r *mediaRepo) Create(ctx context.Context, media *models.Media) error {
	query := `
		INSERT INTO media (
			owner_id,
			object_key,
			original_name,
			mime_type,
			size,
			storage_provider
		)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at;
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		media.OwnerID,
		media.ObjectKey,
		media.OriginalName,
		media.MimeType,
		media.Size,
		media.StorageProvider,
	).Scan(
		&media.ID,
		&media.CreatedAt,
	)
	if err != nil {
		return apperror.Wrap("repository", "CreateMedia", "failed to create media", err)
	}

	return nil
}

func (r *mediaRepo) GetByID(ctx context.Context, id string) (*models.Media, error) {
	query := `
		SELECT
			id,
			owner_id,
			object_key,
			original_name,
			mime_type,
			size,
			storage_provider,
			created_at
		FROM media
		WHERE id = $1;
	`

	var media models.Media

	err := r.db.QueryRowContext(
		ctx,
		query,
		id,
	).Scan(
		&media.ID,
		&media.OwnerID,
		&media.ObjectKey,
		&media.OriginalName,
		&media.MimeType,
		&media.Size,
		&media.StorageProvider,
		&media.CreatedAt,
	)
	if err != nil {
		return nil, apperror.Wrap("repository", "GetMediaByID", "failed to get media by id", err)
	}

	return &media, nil
}

func (r *mediaRepo) GetMany(ctx context.Context, ids []string) ([]models.Media, error) {
	query := `
		SELECT
			id,
			owner_id,
			object_key,
			original_name,
			mime_type,
			size,
			storage_provider,
			created_at
		FROM media
		WHERE id = ANY($1);
	`

	rows, err := r.db.QueryContext(
		ctx,
		query,
		pq.Array(ids),
	)
	if err != nil {
		return nil, apperror.Wrap("repository", "GetMedias", "failed to execute query", err)
	}
	defer rows.Close()

	var medias []models.Media
	for rows.Next() {
		var media models.Media
		err = rows.Scan(
			&media.ID,
			&media.OwnerID,
			&media.ObjectKey,
			&media.OriginalName,
			&media.MimeType,
			&media.Size,
			&media.StorageProvider,
			&media.CreatedAt,
		)
		if err != nil {
			return nil, apperror.Wrap("repository", "GetMedias", "failed to scan media", err)
		}

		medias = append(medias, media)
	}

	if err = rows.Err(); err != nil {
		return nil, apperror.Wrap("repository", "GetMedias", "failed to check rows to error", err)
	}

	return medias, nil
}

func (r *mediaRepo) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM media
		WHERE id = $1;
	`

	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return apperror.Wrap("repository", "DeleteMedia", "failed to delete media", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "DeleteMedia", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "DeleteMedia", "no rows effected on media delete", err)
	}

	return nil
}
