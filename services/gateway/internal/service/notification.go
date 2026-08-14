package service

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/gen/go/common"
	"github.com/diyorbeknematov/minitwitter/gen/go/notification"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/mapper"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"github.com/google/uuid"
)

type notificationService struct {
	client *grpcclient.Client
}

func NewNotificationService(c *grpcclient.Client) *notificationService {
	return &notificationService{
		client: c,
	}
}

func (s *notificationService) GetNotification(
	ctx context.Context,
	id uuid.UUID,
) (dto.Notification, error) {
	resp, err := s.client.Notification.GetNotification(
		ctx,
		&notification.GetNotificationRequest{
			Id: id.String(),
		},
	)
	if err != nil {
		return dto.Notification{}, apperror.Wrap(
			"service",
			"GetNotification",
			"failed to get notification",
			err,
		)
	}

	notif, err := mapper.ToNotification(resp)
	if err != nil {
		return dto.Notification{}, apperror.Wrap(
			"service",
			"GetNotification",
			"failed to map notification",
			err,
		)
	}

	return notif, nil
}

func (s *notificationService) GetNotifications(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int32,
) (dto.GetNotificationsResp, error) {

	resp, err := s.client.Notification.GetNotifications(
		ctx,
		&notification.GetNotificationsRequest{
			UserId: userID.String(),
			Pagination: &common.PaginationRequest{
				Page:  page,
				Limit: limit,
			},
		},
	)
	if err != nil {
		return dto.GetNotificationsResp{}, apperror.Wrap(
			"service",
			"GetNotifications",
			"failed to get notifications",
			err,
		)
	}

	dtoResp, err := mapper.ToGetNotificationsResp(resp)
	if err != nil {
		return dto.GetNotificationsResp{}, apperror.Wrap(
			"service",
			"GetNotifications",
			"failed to map notifications",
			err,
		)
	}

	return dtoResp, nil
}

func (s *notificationService) GetUnreadCount(
	ctx context.Context,
	userID uuid.UUID,
) (dto.GetUnreadCountResp, error) {

	resp, err := s.client.Notification.GetUnreadCount(
		ctx,
		&notification.GetUnreadCountRequest{
			UserId: userID.String(),
		},
	)
	if err != nil {
		return dto.GetUnreadCountResp{}, apperror.Wrap(
			"service",
			"GetUnreadCount",
			"failed to get unread count",
			err,
		)
	}

	return mapper.ToGetUnreadCountResp(resp), nil
}

func (s *notificationService) MarkAsRead(
	ctx context.Context,
	id uuid.UUID,
) error {

	_, err := s.client.Notification.MarkAsRead(
		ctx,
		&notification.MarkAsReadRequest{
			Id: id.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"MarkAsRead",
			"failed to mark notification as read",
			err,
		)
	}

	return nil
}

func (s *notificationService) MarkAllAsRead(
	ctx context.Context,
	userID uuid.UUID,
) error {

	_, err := s.client.Notification.MarkAllAsRead(
		ctx,
		&notification.MarkAllAsReadRequest{
			UserId: userID.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"MarkAllAsRead",
			"failed to mark all notifications as read",
			err,
		)
	}

	return nil
}

func (s *notificationService) DeleteNotification(
	ctx context.Context,
	id uuid.UUID,
) error {

	_, err := s.client.Notification.DeleteNotification(
		ctx,
		&notification.DeleteNotificationRequest{
			Id: id.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"DeleteNotification",
			"failed to delete notification",
			err,
		)
	}

	return nil
}
