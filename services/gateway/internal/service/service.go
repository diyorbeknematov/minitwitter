package service

import "github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"

type Service struct {
	Auth         *authService
	User         *userService
	Tweet        *tweetService
	Media        *mediaService
	Notification *notificationService
}

func New(clients *grpcclient.Client) *Service {
	return &Service{
		Auth:         NewAuthService(clients),
		User:         NewUserService(clients),
		Tweet:        NewTweetService(clients),
		Media:        NewMediaService(clients),
		Notification: NewNotificationService(clients),
	}
}
