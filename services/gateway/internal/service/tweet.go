package service

import "github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"

type tweetService struct {
	client *grpcclient.Client
}

func NewTweetService(c *grpcclient.Client) *tweetService {
	return &tweetService{
		client: c,
	}
}
