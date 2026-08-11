package service

import "github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"

type notificationService struct {
	clinet *grpcclient.Client
}

func NewNotificationService(c *grpcclient.Client) *notificationService {
	return &notificationService{
		clinet: c,
	}
}
