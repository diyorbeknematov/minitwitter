package grpcclient

import (
	"errors"

	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/gen/go/notification"
	"github.com/diyorbeknematov/minitwitter/gen/go/tweet"
	"github.com/diyorbeknematov/minitwitter/gen/go/user"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/config"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"google.golang.org/grpc"
)

type Client struct {
	User         user.UserServiceClient
	Tweet        tweet.TweetServiceClient
	Media        media.MediaServiceClient
	Notification notification.NotificationServiceClient

	userConn         *grpc.ClientConn
	tweetConn        *grpc.ClientConn
	mediaConn        *grpc.ClientConn
	notificationConn *grpc.ClientConn
}

func New(cfg *config.Config) (*Client, error) {
	userConn, err := newConn(cfg.GRPC.User)
	if err != nil {
		return nil, apperror.Wrap("grpcclient", "NewClient", "failed to create connection to gRPC user", err)
	}

	tweetConn, err := newConn(cfg.GRPC.Tweet)
	if err != nil {
		return nil, apperror.Wrap("grpcclient", "NewClient", "failed to create connection to gRPC tweet", err)
	}

	mediaConn, err := newConn(cfg.GRPC.Media)
	if err != nil {
		return nil, apperror.Wrap("grpcclient", "NewClient", "failed to create connection to gRPC media", err)
	}

	notificationConn, err := newConn(cfg.GRPC.Notification)
	if err != nil {
		return nil, apperror.Wrap("grpcclient", "NewClient", "failed to create connection to gRPC notification", err)
	}

	return &Client{
		userConn: userConn,
		User:     user.NewUserServiceClient(userConn),

		tweetConn: tweetConn,
		Tweet:     tweet.NewTweetServiceClient(tweetConn),

		mediaConn: mediaConn,
		Media:     media.NewMediaServiceClient(mediaConn),

		notificationConn: notificationConn,
		Notification:     notification.NewNotificationServiceClient(notificationConn),
	}, nil
}

func (c *Client) Close() error {
	var errs []error

	for _, conn := range []*grpc.ClientConn{
		c.userConn,
		c.tweetConn,
		c.mediaConn,
		c.notificationConn,
	} {
		if conn == nil {
			continue
		}

		if err := conn.Close(); err != nil {
			errs = append(errs, err)
		}
	}

	return errors.Join(errs...)
}
