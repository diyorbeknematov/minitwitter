package service

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/gen/go/common"
	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/gen/go/tweet"
	"github.com/diyorbeknematov/minitwitter/gen/go/user"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/mapper"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"github.com/google/uuid"
)

type tweetService struct {
	client *grpcclient.Client
}

func NewTweetService(c *grpcclient.Client) *tweetService {
	return &tweetService{
		client: c,
	}
}

func (s *tweetService) CreateTweet(
	ctx context.Context,
	userID uuid.UUID,
	req dto.CreateTweetReq,
) (dto.Tweet, error) {

	if len(req.MediaIDs) > 0 {

		medias, err := s.client.Media.GetMedias(
			ctx,
			mapper.ToGetMediasReq(req.MediaIDs),
		)
		if err != nil {
			return dto.Tweet{}, err
		}

		for _, media := range medias.Medias {

			ownerID, _ := uuid.Parse(media.OwnerId)

			if ownerID != userID {
				return dto.Tweet{}, apperror.Wrap(
					"service", "CreateTweet", "failed to cheek media ids", err,
				)
			}
		}
	}

	resp, err := s.client.Tweet.CreateTweet(
		ctx,
		mapper.ToCreateTweetReq(userID, req),
	)
	if err != nil {
		return dto.Tweet{}, err
	}

	return mapper.ToTweet(resp.Tweet)
}

func (s *tweetService) UpdateTweet(ctx context.Context, tweetID uuid.UUID, req dto.UpdateTweetReq) error {
	_, err := s.client.Tweet.UpdateTweet(ctx, mapper.ToUpdateTweetReq(tweetID, req))
	if err != nil {
		return apperror.Wrap(
			"service", "UpdateTweet", "failed to update tweet", err,
		)
	}

	return nil
}

func (s *tweetService) DeleteTweet(ctx context.Context, tweetID uuid.UUID) error {
	_, err := s.client.Tweet.DeleteTweet(ctx, &tweet.DeleteTweetRequest{TweetId: tweetID.String()})
	if err != nil {
		return apperror.Wrap(
			"service", "DeleteTweet", "failed to delete tweet", err,
		)
	}

	return nil
}

func (s *tweetService) GetTweet(ctx context.Context, tweetID uuid.UUID) (dto.Tweet, error) {
	resp, err := s.client.Tweet.GetTweet(ctx, &tweet.GetTweetRequest{TweetId: tweetID.String()})
	if err != nil {
		return dto.Tweet{}, apperror.Wrap(
			"service", "GetTweet", "failed to get tweet by id", err,
		)
	}

	twt, err := mapper.ToTweet(resp)
	if err != nil {
		return dto.Tweet{}, apperror.Wrap(
			"service", "GetTweet", "failed to parse to dto.Tweet", err,
		)
	}

	usrResp, err := s.client.User.GetUserById(ctx, &user.GetUserByIdRequest{Id: twt.AuthorID.String()})
	if err != nil {
		return dto.Tweet{}, apperror.Wrap(
			"service", "GetTweet", "failed to get author", err,
		)
	}

	u, err := mapper.ToUser(usrResp)
	if err != nil {
		return dto.Tweet{}, apperror.Wrap(
			"service", "GetTweet", "failed to parse dto.User", err,
		)
	}

	mediaSet := make(map[string]struct{})
	for _, m := range twt.MediaIDs {
		mediaSet[m.String()] = struct{}{}
	}

	if u.AvatarMediaID != uuid.Nil {
		mediaSet[u.AvatarMediaID.String()] = struct{}{}
	}

	mediaIDs := make([]string, 0, len(mediaSet))
	for id := range mediaSet {
		mediaIDs = append(mediaIDs, id)
	}

	// Medias
	mediaResp, err := s.client.Media.GetMedias(
		ctx,
		&media.GetMediasRequest{
			MediaIds: mediaIDs,
		},
	)
	if err != nil {
		return dto.Tweet{}, apperror.Wrap(
			"service", "GetTweet", "failed to get medias", err,
		)
	}

	medias, err := mapper.ToMediaMap(mediaResp)
	if err != nil {
		return dto.Tweet{}, apperror.Wrap(
			"service", "GetTweet", "failed to map medias", err,
		)
	}

	if m, ok := medias[u.AvatarMediaID.String()]; ok {
		u.AvatarURL = m.Url
		twt.Author = u
	}

	for _, mediaID := range twt.MediaIDs {
		if media, ok := medias[mediaID.String()]; ok {
			twt.MediaURLs = append(twt.MediaURLs, media.Url)
		}
	}

	return twt, nil
}

// GetTweetsByUser
func (s *tweetService) GetTweetsByUser(ctx context.Context, userID uuid.UUID, req dto.GetTweetByUserReq) (dto.GetTweetByUserResp, error) {
	resp, err := s.client.Tweet.GetTweetsByUser(ctx, mapper.ToGetTweetsByUserReq(userID, req))
	if err != nil {
		return dto.GetTweetByUserResp{}, apperror.Wrap(
			"service", "GetTweetsByUser", "failed to get tweets by user", err,
		)
	}

	twts, err := mapper.ToGetTweetsByUserResp(resp)
	if err != nil {
		return dto.GetTweetByUserResp{}, apperror.Wrap(
			"service", "GetTweetsByUser", "failed to parse to dto.GetTweetByUserResp", err,
		)
	}

	if len(twts.Tweets) == 0 {
		return twts, nil
	}

	// Bitta user — barcha tweet shu userniki
	usrResp, err := s.client.User.GetUserById(ctx, &user.GetUserByIdRequest{Id: userID.String()})
	if err != nil {
		return dto.GetTweetByUserResp{}, apperror.Wrap(
			"service", "GetTweetsByUser", "failed to get user", err,
		)
	}

	u, err := mapper.ToUser(usrResp)
	if err != nil {
		return dto.GetTweetByUserResp{}, apperror.Wrap(
			"service", "GetTweetsByUser", "failed to parse dto.User", err,
		)
	}

	// Barcha tweet media ID lar + avatar media ID — bitta setga
	mediaSet := make(map[string]struct{})
	for _, t := range twts.Tweets {
		for _, m := range t.MediaIDs {
			mediaSet[m.String()] = struct{}{}
		}
	}
	if u.AvatarMediaID != uuid.Nil {
		mediaSet[u.AvatarMediaID.String()] = struct{}{}
	}

	if len(mediaSet) > 0 {
		mediaIDs := make([]string, 0, len(mediaSet))
		for id := range mediaSet {
			mediaIDs = append(mediaIDs, id)
		}

		mediaResp, err := s.client.Media.GetMedias(ctx, &media.GetMediasRequest{MediaIds: mediaIDs})
		if err != nil {
			return dto.GetTweetByUserResp{}, apperror.Wrap(
				"service", "GetTweetsByUser", "failed to get medias", err,
			)
		}

		medias, err := mapper.ToMediaMap(mediaResp)
		if err != nil {
			return dto.GetTweetByUserResp{}, apperror.Wrap(
				"service", "GetTweetsByUser", "failed to map medias", err,
			)
		}

		if m, ok := medias[u.AvatarMediaID.String()]; ok {
			u.AvatarURL = m.Url
		}

		for i := range twts.Tweets {
			for _, mediaID := range twts.Tweets[i].MediaIDs {
				if m, ok := medias[mediaID.String()]; ok {
					twts.Tweets[i].MediaURLs = append(twts.Tweets[i].MediaURLs, m.Url)
				}
			}
		}
	}

	// Bitta userni barcha tweet'larga tarqatamiz
	for i := range twts.Tweets {
		twts.Tweets[i].Author = u
	}

	return twts, nil
}

func (s *tweetService) GetTimeline(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int32,
) (dto.GetTimelineResp, error) {

	// Following user IDs
	following, err := s.client.User.GetFollowingIds(
		ctx,
		&user.GetFollowingIdsRequest{
			Id: userID.String(),
		},
	)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to get following ids", err,
		)
	}

	// Timeline
	resp, err := s.client.Tweet.GetTimeline(
		ctx,
		&tweet.GetTimelineRequest{
			UserIds: following.Ids,
			Pagination: &common.PaginationRequest{
				Page:  page,
				Limit: limit,
			},
		},
	)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to get timeline", err,
		)
	}

	timeline, err := mapper.ToGetTimelineResp(resp)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to map timeline", err,
		)
	}

	if len(timeline.Tweets) == 0 {
		return timeline, nil
	}

	// Collect unique IDs
	userSet := make(map[string]struct{})
	mediaSet := make(map[string]struct{})

	for _, t := range timeline.Tweets {
		userSet[t.AuthorID.String()] = struct{}{}

		for _, mediaID := range t.MediaIDs {
			mediaSet[mediaID.String()] = struct{}{}
		}
	}

	userIDs := make([]string, 0, len(userSet))
	for id := range userSet {
		userIDs = append(userIDs, id)
	}

	// Users
	userResp, err := s.client.User.GetUsersByIds(
		ctx,
		&user.GetUsersByIdsRequest{
			UserIds: userIDs,
		},
	)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to get users", err,
		)
	}

	users, err := mapper.ToUsersMap(userResp)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to map users", err,
		)
	}

	for _, u := range users {
		if u.AvatarMediaID != uuid.Nil {
			mediaSet[u.AvatarMediaID.String()] = struct{}{}
		}
	}

	mediaIDs := make([]string, 0, len(mediaSet))
	for id := range mediaSet {
		mediaIDs = append(mediaIDs, id)
	}

	// Medias
	mediaResp, err := s.client.Media.GetMedias(
		ctx,
		&media.GetMediasRequest{
			MediaIds: mediaIDs,
		},
	)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to get medias", err,
		)
	}

	medias, err := mapper.ToMediaMap(mediaResp)
	if err != nil {
		return dto.GetTimelineResp{}, apperror.Wrap(
			"service", "GetTimeline", "failed to map medias", err,
		)
	}

	// Userlarga avatar URL/media ni to'ldiramiz
	for id, u := range users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}
		if m, ok := medias[u.AvatarMediaID.String()]; ok {
			u.AvatarURL = m.Url
			users[id] = u
		}
	}

	for i := range timeline.Tweets {

		t := &timeline.Tweets[i]

		if usr, ok := users[t.AuthorID.String()]; ok {
			t.Author = usr
		}

		t.MediaURLs = make([]string, 0, len(t.MediaIDs))

		for _, mediaID := range t.MediaIDs {

			if media, ok := medias[mediaID.String()]; ok {
				t.MediaURLs = append(t.MediaURLs, media.Url)
			}
		}
	}

	return timeline, nil
}

func (s *tweetService) LikeTweet(
	ctx context.Context,
	userID, tweetID uuid.UUID,
) error {
	_, err := s.client.Tweet.LikeTweet(
		ctx,
		&tweet.LikeTweetRequest{
			TweetId: tweetID.String(),
			UserId:  userID.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"LikeTweet",
			"failed to like tweet",
			err,
		)
	}

	return nil
}

func (s *tweetService) UnlikeTweet(
	ctx context.Context,
	userID, tweetID uuid.UUID,
) error {
	_, err := s.client.Tweet.UnlikeTweet(
		ctx,
		&tweet.UnlikeTweetRequest{
			TweetId: tweetID.String(),
			UserId:  userID.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"UnlikeTweet",
			"failed to unlike tweet",
			err,
		)
	}

	return nil
}

func (s *tweetService) Retweet(
	ctx context.Context,
	userID, tweetID uuid.UUID,
) error {
	_, err := s.client.Tweet.Retweet(
		ctx,
		&tweet.RetweetRequest{
			TweetId: tweetID.String(),
			UserId:  userID.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"Retweet",
			"failed to retweet",
			err,
		)
	}

	return nil
}

func (s *tweetService) UndoRetweet(
	ctx context.Context,
	userID, tweetID uuid.UUID,
) error {
	_, err := s.client.Tweet.UndoRetweet(
		ctx,
		&tweet.UndoRetweetRequest{
			TweetId: tweetID.String(),
			UserId:  userID.String(),
		},
	)
	if err != nil {
		return apperror.Wrap(
			"service",
			"UndoRetweet",
			"failed to undo retweet",
			err,
		)
	}

	return nil
}
