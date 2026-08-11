package mapper

import (
	"errors"

	"github.com/diyorbeknematov/minitwitter/gen/go/common"
	"github.com/diyorbeknematov/minitwitter/gen/go/notification"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/google/uuid"
)

func ToNotification(pb *notification.Notification) (dto.Notification, error) {
	id, err := uuid.Parse(pb.Id)
	if err != nil {
		return dto.Notification{}, err
	}

	actorID, err := uuid.Parse(pb.ActorId)
	if err != nil {
		return dto.Notification{}, err
	}

	entityID, err := uuid.Parse(pb.EntityId)
	if err != nil {
		return dto.Notification{}, err
	}

	var notificationType dto.NotificationType

	switch pb.Type {
	case notification.NotificationType_FOLLOW:
		notificationType = dto.NotificationTypeFollow

	case notification.NotificationType_LIKE:
		notificationType = dto.NotificationTypeLike

	case notification.NotificationType_RETWEET:
		notificationType = dto.NotificationTypeRetweet

	case notification.NotificationType_REPLY:
		notificationType = dto.NotificationTypeReply

	default:
		return dto.Notification{}, errors.New("unknown notification type")
	}

	return dto.Notification{
		ID:        id,
		ActorID:   actorID,
		Type:      notificationType,
		EntityID:  entityID,
		IsRead:    pb.IsRead,
		CreatedAt: pb.CreatedAt.AsTime(),
	}, nil
}

func ToGetNotificationsResp(
	pb *notification.GetNotificationsResponse,
) (dto.GetNotificationsResp, error) {

	notifications := make([]dto.Notification, len(pb.Notifications))

	for i, n := range pb.Notifications {
		notificationDTO, err := ToNotification(n)
		if err != nil {
			return dto.GetNotificationsResp{}, err
		}

		notifications[i] = notificationDTO
	}

	return dto.GetNotificationsResp{
		Notifications: notifications,
		Pagination: dto.Pagination{
			Page:  pb.Pagination.Page,
			Limit: pb.Pagination.Limit,
			Total: pb.Pagination.Total,
		},
	}, nil
}

func ToGetNotificationsRequest(
	userID uuid.UUID,
	req dto.GetNotificationsReq,
) *notification.GetNotificationsRequest {

	return &notification.GetNotificationsRequest{
		UserId: userID.String(),
		Pagination: &common.PaginationRequest{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}
}

func ToGetUnreadCountResponse(
	pb *notification.GetUnreadCountResponse,
) dto.GetUnreadCountResp {

	return dto.GetUnreadCountResp{
		Count: pb.Count,
	}
}