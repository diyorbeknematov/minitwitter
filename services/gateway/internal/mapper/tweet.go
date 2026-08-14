package mapper

import (
	"github.com/diyorbeknematov/minitwitter/gen/go/common"
	"github.com/diyorbeknematov/minitwitter/gen/go/tweet"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/google/uuid"
)

func ToTweet(pb *tweet.Tweet) (dto.Tweet, error) {
	id, err := uuid.Parse(pb.Id)
	if err != nil {
		return dto.Tweet{}, err
	}

	authorID, err := uuid.Parse(pb.AuthorId)
	if err != nil {
		return dto.Tweet{}, err
	}

	var replyID *uuid.UUID

	if pb.ReplyToTweetId != "" {
		id, err := uuid.Parse(pb.ReplyToTweetId)
		if err != nil {
			return dto.Tweet{}, err
		}

		replyID = &id
	}

	mediaIDs := make([]uuid.UUID, len(pb.MediaIds))
	for i, mID := range pb.MediaIds {
		id, err := uuid.Parse(mID)
		if err != nil {
			return dto.Tweet{}, err
		}
		mediaIDs[i] = id
	}

	return dto.Tweet{
		ID:             id,
		AuthorID:       authorID,
		Content:        pb.Content,
		MediaIDs:       mediaIDs,
		ReplyToTweetID: replyID,
		LikesCount:     pb.LikesCount,
		RetweetsCount:  pb.RetweetsCount,
		CreatedAt:      pb.CreatedAt.AsTime(),
		UpdatedAt:      pb.UpdatedAt.AsTime(),
	}, nil
}

func ToCreateTweetReq(
	userID uuid.UUID,
	req dto.CreateTweetReq,
) *tweet.CreateTweetRequest {
	mediaIDs := make([]string, len(req.MediaIDs))
	for i, mID := range req.MediaIDs {
		mediaIDs[i] = mID.String()
	}

	var replyID string

	if req.ReplyToTweetID != nil {
		replyID = req.ReplyToTweetID.String()
	}

	return &tweet.CreateTweetRequest{
		AuthorId:       userID.String(),
		Content:        req.Content,
		MediaIds:       mediaIDs,
		ReplyToTweetId: replyID,
	}
}

func ToUpdateTweetReq(
	tweetID uuid.UUID,
	req dto.UpdateTweetReq,
) *tweet.UpdateTweetRequest {
	mediaIDs := make([]string, len(req.MediaIDs))
	for i, mID := range req.MediaIDs {
		mediaIDs[i] = mID.String()
	}

	return &tweet.UpdateTweetRequest{
		TweetId:  tweetID.String(),
		Content:  req.Content,
		MediaIds: mediaIDs,
	}
}

func ToGetTweetsByUserReq(userID uuid.UUID, req dto.GetTweetByUserReq) *tweet.GetTweetsByUserRequest {
	return &tweet.GetTweetsByUserRequest{
		UserId: userID.String(),
		Pagination: &common.PaginationRequest{
			Page:  req.Page,
			Limit: req.Limit,
		},
	}
}

func ToGetTweetsByUserResp(
	pb *tweet.GetTweetsByUserResponse,
) (dto.GetTweetByUserResp, error) {

	tweets := make([]dto.Tweet, len(pb.Tweets))

	for i, t := range pb.Tweets {
		tweetDTO, err := ToTweet(t)
		if err != nil {
			return dto.GetTweetByUserResp{}, err
		}

		tweets[i] = tweetDTO
	}

	return dto.GetTweetByUserResp{
		Tweets: tweets,
		Pagination: dto.Pagination{
			Page:  pb.Pagination.Page,
			Limit: pb.Pagination.Limit,
			Total: pb.Pagination.Total,
		},
	}, nil
}

func ToGetTimelineResp(
	pb *tweet.GetTimelineResponse,
) (dto.GetTimelineResp, error) {

	tweets := make([]dto.Tweet, len(pb.Tweets))

	for i, t := range pb.Tweets {
		tweetDTO, err := ToTweet(t)
		if err != nil {
			return dto.GetTimelineResp{}, err
		}

		tweets[i] = tweetDTO
	}

	return dto.GetTimelineResp{
		Tweets: tweets,
		Pagination: dto.Pagination{
			Page:  pb.Pagination.Page,
			Limit: pb.Pagination.Limit,
			Total: pb.Pagination.Total,
		},
	}, nil
}
