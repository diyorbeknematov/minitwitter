package service

import "github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"

type mediaService struct {
	client *grpcclient.Client
}

func NewMediaService(c *grpcclient.Client) *mediaService {
	return &mediaService{
		client: c,
	}
}
