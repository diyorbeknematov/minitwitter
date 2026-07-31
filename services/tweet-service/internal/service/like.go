package service

import (
	"context"
	"time"

	"github.com/diyorbeknematov/minitwitter/gen/go/tweet"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *tweetService) LikeTweet(ctx context.Context, req *tweet.LikeTweetRequest) (*emptypb.Empty, error) {
	twtID, err := uuid.Parse(req.TweetId)
	if err != nil {
		return nil, apperror.Wrap("service", "LikeTweet", "failed to parse tweeetId string to uuid", err)
	}

	usrId, err := uuid.Parse(req.UserId)
	if err != nil {
		return nil, apperror.Wrap("service", "LikeTweet", "failed to parse user id from string to uuid", err)
	}

	err = s.repo.Like.Create(ctx, &models.Like{
		TweetID:   twtID,
		UserID:    usrId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, apperror.Wrap("service", "LikeTweet", "failed to like tweet", err)
	}

	return nil, nil
}

func (s *tweetService) UnlikeTweet(ctx context.Context, req *tweet.UnlikeTweetRequest) (*emptypb.Empty, error) {
	err := s.repo.Like.Delete(ctx, req.TweetId, req.UserId)
	if err != nil {
		return nil, apperror.Wrap("service", "UnlikeTweet", "failed to unlike tweet", err)
	}

	return nil, nil
}
