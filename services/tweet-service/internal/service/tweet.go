package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/diyorbeknematov/minitwitter/gen/go/common"
	"github.com/diyorbeknematov/minitwitter/gen/go/tweet"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/models"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/internal/repository"
	"github.com/diyorbeknematov/minitwitter/services/tweet-service/pkg/apperror"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type tweetService struct {
	repo   *repository.Repository
	logger *slog.Logger

	tweet.UnimplementedTweetServiceServer
}

func NewTweetService(repo *repository.Repository, logger *slog.Logger) *tweetService {
	return &tweetService{
		repo:   repo,
		logger: logger,
	}
}

func (s *tweetService) CreateTweet(ctx context.Context, req *tweet.CreateTweetRequest) (*tweet.CreateTweetResponse, error) {
	authorID, err := uuid.Parse(req.GetAuthorId())
	if err != nil {
		return nil, apperror.Wrap("service", "CreateTweet", "failed to parse author id to uuid", err)
	}

	replyID, err := uuid.Parse(req.GetReplyToTweetId())
	if err != nil {
		return nil, apperror.Wrap("service", "CreateTweet", "failed to parse replytotweetid to uuid", err)
	}

	twt := &models.Tweet{
		AuthorID:       authorID,
		Content:        req.GetContent(),
		ReplyToTweetID: &replyID,
		UpdatedAt:      time.Now(),
	}

	err = s.repo.Tweet.Create(ctx, twt)
	if err != nil {
		return nil, apperror.Wrap("service", "CreateTweet", "failed to create tweet", err)
	}

	return &tweet.CreateTweetResponse{
		Tweet: &tweet.Tweet{
			Id:             twt.ID.String(),
			AuthorId:       twt.AuthorID.String(),
			Content:        twt.Content,
			ReplyToTweetId: twt.ReplyToTweetID.String(),
			CreatedAt:      timestamppb.New(twt.CreatedAt),
		},
	}, nil
}

func (s *tweetService) UpdateTweet(ctx context.Context, req *tweet.UpdateTweetRequest) (*emptypb.Empty, error) {
	twtID, err := uuid.Parse(req.TweetId)
	if err != nil {
		return nil, apperror.Wrap("service", "UpdateTweet", "failed to parse twtID to uuid", err)
	}

	err = s.repo.Tweet.Update(ctx, models.Tweet{
		ID:      twtID,
		Content: req.Content,
	})
	if err != nil {
		return nil, apperror.Wrap("service", "UpdateTweet", "failed to update tweet", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *tweetService) DeleteTweet(ctx context.Context, req *tweet.DeleteTweetRequest) (*emptypb.Empty, error) {
	err := s.repo.Tweet.Delete(ctx, req.TweetId)
	if err != nil {
		return nil, apperror.Wrap("service", "DeleteTweet", "failed to delete tweet", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *tweetService) GetTweet(ctx context.Context, req *tweet.GetTweetRequest) (*tweet.Tweet, error) {
	twt, err := s.repo.Tweet.GetByID(ctx, req.TweetId)
	if err != nil {
		return nil, apperror.Wrap("service", "GetByID", "failed to get tweet by id", err)
	}

	return &tweet.Tweet{
		Id:             twt.ID.String(),
		AuthorId:       twt.AuthorID.String(),
		Content:        twt.Content,
		ReplyToTweetId: twt.ReplyToTweetID.String(),
		CreatedAt:      timestamppb.New(twt.CreatedAt),
		UpdatedAt:      timestamppb.New(twt.UpdatedAt),
	}, nil
}

func (s *tweetService) GetTweetsByUser(ctx context.Context, req *tweet.GetTweetsByUserRequest) (*tweet.GetTweetsByUserResponse, error) {
	twts, total, err := s.repo.Tweet.GetByUser(ctx, req.UserId, req.Pagination.Limit, (req.Pagination.Page-1)*req.Pagination.Limit)
	if err != nil {
		return nil, apperror.Wrap("service", "GetTweetsByUser", "failed to get tweets by username", err)
	}

	return &tweet.GetTweetsByUserResponse{
		Tweets: s.toProtoTweets(twts),
		Pagination: &common.PaginationResponse{
			Page:  req.Pagination.Page,
			Limit: req.Pagination.Limit,
			Total: int64(total),
		},
	}, nil
}

func (s *tweetService) GetTimeline(ctx context.Context, req *tweet.GetTimelineRequest) (*tweet.GetTimelineResponse, error) {
	uIds := make([]uuid.UUID, 0, len(req.UserIds))

	for _, id := range req.UserIds {
		uId, err := uuid.Parse(id)
		if err != nil {
			return nil, apperror.Wrap("serivice", "GetTimeline", "failed to parse user ids to uuid", err)
		}

		uIds = append(uIds, uId)
	}

	twts, total, err := s.repo.Tweet.GetTimeline(
		ctx,
		uIds,
		req.Pagination.Limit,
		(req.Pagination.Page-1)*req.Pagination.Limit,
	)
	if err != nil {
		return nil, apperror.Wrap("service", "GetTimeline", "failed to get timeline", err)
	}

	return &tweet.GetTimelineResponse{
		Tweets: s.toProtoTweets(twts),
		Pagination: &common.PaginationResponse{
			Page:  req.Pagination.Page,
			Limit: req.Pagination.Limit,
			Total: int64(total),
		},
	}, nil
}

func (s *tweetService) Retweet(ctx context.Context, req *tweet.RetweetRequest) (*emptypb.Empty, error) {
	twtID, err := uuid.Parse(req.TweetId)
	if err != nil {
		return nil, apperror.Wrap("service", "Retweet", "failed to parse twtID from string to uuid", err)
	}

	usrId, err := uuid.Parse(req.TweetId)
	if err != nil {
		return nil, apperror.Wrap("service", "Retweet", "failed to parse user id from string to uuid", err)
	}

	err = s.repo.Tweet.CreateRetweet(ctx, models.Retweet{
		TweetID:   twtID,
		UserID:    usrId,
		CreatedAt: time.Now(),
	})
	if err != nil {
		return nil, apperror.Wrap("service", "Retweet", "failed to retweet tweet", err)
	}

	return &emptypb.Empty{}, nil
}

func (s *tweetService) UndoRetweet(ctx context.Context, req *tweet.UndoRetweetRequest) (*emptypb.Empty, error) {
	err := s.repo.Tweet.DeleteRetweet(ctx, req.TweetId, req.UserId)
	if err != nil {
		return nil, apperror.Wrap("service", "UndoRetweet", "failed to undo retweet", err)
	}

	return nil, nil
}

func (s *tweetService) toProtoTweet(u models.Tweet) *tweet.Tweet {
	return &tweet.Tweet{
		Id:             u.ID.String(),
		AuthorId:       u.AuthorID.String(),
		Content:        u.Content,
		ReplyToTweetId: u.ReplyToTweetID.String(),
	}
}

func (s *tweetService) toProtoTweets(twts []models.Tweet) []*tweet.Tweet {
	result := make([]*tweet.Tweet, 0, len(twts))

	for _, u := range twts {
		result = append(result, s.toProtoTweet(u))
	}

	return result
}
