package postgres

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/pkg/apperror"
	"github.com/jmoiron/sqlx"
)

type notificationRepo struct {
	db *sqlx.DB
}

func NewNotificationRepo(db *sqlx.DB) *notificationRepo {
	return &notificationRepo{
		db: db,
	}
}

func (r *notificationRepo) Create(ctx context.Context, notif *models.Notification) error {
	query := `
		INSERT INTO notifications (
			user_id,
			actor_id,
			type,
			entity_id,
			is_read
		)
		VALUES($1, $2, $3, $4, $5)
		RETURNING id, created_at;
	`

	err := r.db.QueryRowContext(
		ctx,
		query,
		notif.UserID,
		notif.ActorID,
		notif.Type,
		notif.EntityID,
		notif.IsRead,
	).Scan(
		&notif.ID,
		&notif.CreatedAt,
	)
	if err != nil {
		return apperror.Wrap("repository", "CreateNotification", "failed to create notification", err)
	}

	return nil
}
func (r *notificationRepo) GetByID(ctx context.Context, id string) (*models.Notification, error) {
	query := `
		SELECT
			id,
			user_id,
			actor_id,
			type,
			entity_id,
			is_read,
			created_at
		FROM notifications
		WHERE id = $1;	
	`
	var notif models.Notification
	err := r.db.QueryRow(query, id).Scan(
		&notif.ID,
		&notif.UserID,
		&notif.ActorID,
		&notif.Type,
		&notif.EntityID,
		&notif.IsRead,
		&notif.CreatedAt,
	)
	if err != nil {
		return nil, apperror.Wrap("repository", "GetNotificationByID", "failed to get notification by id", err)
	}

	return &notif, nil
}
func (r *notificationRepo) GetMany(ctx context.Context, req models.GetNotifications) ([]models.Notification, int, error) {
	baseQuery := `
		SELECT
			id,
			user_id,
			actor_id,
			type,
			entity_id,
			is_read,
			created_at
		FROM notifications
		WHERE user_id = $1
		LIMIT $2 OFFSET $3;
	`
	countQuery := `
		SELECT 
			COUNT(*)
		FROM notifications
		WHERE user_id =  $1;
	`

	rows, err := r.db.QueryContext(
		ctx,
		baseQuery,
		req.UserID,
		req.Limit,
		req.Offset,
	)
	if err != nil {
		return nil, 0, apperror.Wrap("repository", "GetNotifications", "failed to execute query", err)
	}
	defer rows.Close()

	var notifs []models.Notification
	for rows.Next() {
		var notif models.Notification
		if err := rows.Scan(
			&notif.ID,
			&notif.UserID,
			&notif.ActorID,
			&notif.Type,
			&notif.EntityID,
			&notif.IsRead,
			&notif.CreatedAt,
		); err != nil {
			return nil, 0, apperror.Wrap("repository", "GetNotifications", "failed to get notifications by user id", err)
		}

		notifs = append(notifs, notif)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, apperror.Wrap("repository", "GetNotifications", "failed to check rows err", err)
	}

	var total int
	if err = r.db.QueryRow(countQuery, req.UserID).Scan(&total); err != nil {
		return nil, 0, apperror.Wrap("repository", "GetNotifications", "failed to get total count", err)
	}

	return notifs, total, nil
}

func (r *notificationRepo) Update(ctx context.Context, id string) error {
	query := `
		UPDATE notifications
		SET
			is_read = true
		WHERE id = $1 and is_read = false;
	`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return apperror.Wrap("repository", "UpdateNotificationByID", "failed to update notification by id", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "UpdateNotificationByID", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "UpdateNotificationByID", "no rows affected to notification update", err)
	}

	return nil
}
func (r *notificationRepo) UpdateAll(ctx context.Context, userID string) error {
	query := `
		UPDATE notifications
		SET
			is_read = true
		WHERE user_id = $1 AND is_read = false;
	`
	res, err := r.db.Exec(query, userID)
	if err != nil {
		return apperror.Wrap("repository", "UpdateAllNotifications", "failed to update all notifications", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "UpdateAllNotifications", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "UpdateAllNotifications", "no rows affected to notifications update", err)
	}

	return nil
}
func (r *notificationRepo) Delete(ctx context.Context, id string) error {
	query := `
		DELETE FROM notifications
		WHERE id = $1;
	`
	res, err := r.db.Exec(query, id)
	if err != nil {
		return apperror.Wrap("repository", "DeleteNotificationByID", "failed to delete notification by id", err)
	}

	rows, err := res.RowsAffected()
	if err != nil {
		return apperror.Wrap("repository", "DeleteNotificationByID", "failed to get rows effected", err)
	}

	if rows == 0 {
		return apperror.Wrap("repository", "DeleteNotificationByID", "no rows affected to notification delete", err)
	}

	return nil
}
func (r *notificationRepo) GetUnreadCount(ctx context.Context, userID string) (int, error) {
	query := `
		SELECT 
			COUNT(*)
		FROM notifications
		WHERE user_id = $1 AND is_read = false;
	`
	var count int
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, apperror.Wrap("repository", "GetUnreadCount", "failed to get unread count", err)
	}
	return count, nil
}
