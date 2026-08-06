package service

import (
	"context"
	"log/slog"

	"github.com/diyorbeknematov/minitwitter/gen/go/common"
	"github.com/diyorbeknematov/minitwitter/gen/go/notification"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/notification-service/pkg/apperror"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type notificationService struct {
	repo   *repository.Repository
	cfg    *config.Config
	logger *slog.Logger

	notification.UnimplementedNotificationServiceServer
}

func NewNotificationService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *notificationService {
	return &notificationService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *notificationService) GetNotification(ctx context.Context, req *notification.GetNotificationRequest) (*notification.Notification, error) {
	ntf, err := s.repo.GetByID(ctx, req.Id)
	if err != nil {
		return nil, apperror.Wrap("service", "GetNotificationByID", "failed to get notification by id", err)
	}

	return s.toProtoNotification(ntf), nil
}

func (s *notificationService) GetNotifications(ctx context.Context, req *notification.GetNotificationsRequest) (*notification.GetNotificationsResponse, error) {
	userID, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, apperror.Wrap("service", "GetNotifications", "failed to parse user id to uuid", err)
	}

	ntfs, total, err := s.repo.GetMany(ctx, models.GetNotifications{
		UserID: userID,
		Limit:  int(req.Pagination.Limit),
		Offset: ((int(req.Pagination.Page) - 1) * int(req.Pagination.Limit)),
	})
	if err != nil {
		return nil, apperror.Wrap("service", "GetNotifications", "failed to get notifications from database", err)
	}

	return &notification.GetNotificationsResponse{
		Notifications: s.toProtoNotifications(ntfs),
		Pagination: &common.PaginationResponse{
			Page:  req.Pagination.Page,
			Limit: req.Pagination.Limit,
			Total: int64(total),
		},
	}, nil
}

func (s *notificationService) GetUnreadCount(ctx context.Context, req *notification.GetUnreadCountRequest) (*notification.GetUnreadCountResponse, error) {
	unreadCount, err := s.repo.GetUnreadCount(ctx, req.UserId)
	if err != nil {
		return nil, apperror.Wrap("service", "GetUnreadCount", "failed to get count from database", err)
	}

	return &notification.GetUnreadCountResponse{
		UserId: req.UserId,
		Count:  int64(unreadCount),
	}, nil
}

func (s *notificationService) MarkAsRead(ctx context.Context, req *notification.MarkAsReadRequest) (*emptypb.Empty, error) {
	err := s.repo.Update(ctx, req.Id)
	if err != nil {
		return nil, apperror.Wrap("service", "MarkAsRead", "failed to update notification from unread to read", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *notificationService) MarkAllAsRead(ctx context.Context, req *notification.MarkAllAsReadRequest) (*emptypb.Empty, error) {
	err := s.repo.UpdateAll(ctx, req.UserId)
	if err != nil {
		return nil, apperror.Wrap("service", "MarkAllAsRead", "failed to update all notifications by user id from unread to read", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *notificationService) DeleteNotification(ctx context.Context, req *notification.DeleteNotificationRequest) (*emptypb.Empty, error) {
	err := s.repo.Delete(ctx, req.Id)
	if err != nil {
		return nil, apperror.Wrap("service", "DeleteNotification", "failed to delete notification", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *notificationService) toProtoNotificationType(t models.NotificationType) notification.NotificationType {
	switch t {
	case models.NotificationFollow:
		return notification.NotificationType_FOLLOW
	case models.NotificationLike:
		return notification.NotificationType_LIKE
	case models.NotificationRetweet:
		return notification.NotificationType_RETWEET
	case models.NotificationReply:
		return notification.NotificationType_REPLY
	default:
		return notification.NotificationType_NOTIFICATION_TYPE_UNSPECIFIED
	}
}

func (s *notificationService) toProtoNotification(ntf *models.Notification) *notification.Notification {
	return &notification.Notification{
		Id:        ntf.ID.String(),
		UserId:    ntf.UserID.String(),
		ActorId:   ntf.ActorID.String(),
		Type:      s.toProtoNotificationType(ntf.Type),
		EntityId:  ntf.EntityID.String(),
		IsRead:    ntf.IsRead,
		CreatedAt: timestamppb.New(ntf.CreatedAt),
	}
}

func (s *notificationService) toProtoNotifications(ntfs []models.Notification) []*notification.Notification {
	result := make([]*notification.Notification, 0, len(ntfs))

	for _, n := range ntfs {
		result = append(result, s.toProtoNotification(&n))
	}

	return result
}
