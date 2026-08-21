package service

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/gen/go/media"
	"github.com/diyorbeknematov/minitwitter/gen/go/user"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/mapper"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"github.com/google/uuid"
)

type userService struct {
	client *grpcclient.Client
}

func NewUserService(c *grpcclient.Client) *userService {
	return &userService{
		client: c,
	}
}

func (s *userService) GetProfile(ctx context.Context, username string) (dto.User, error) {
	resp, err := s.client.User.GetProfile(ctx, &user.GetProfileRequest{Username: username})
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "GetProfile", "failed to get user profile", err)
	}

	dtoUser, err := mapper.ToUser(resp)
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "GetProfile", "failed to parse dto.User", err)
	}

	avatar, err := s.client.Media.GetMedia(ctx, &media.GetMediaRequest{MediaId: resp.AvatarMediaId})
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "GetProfile", "failed to get user avatar", err)
	}

	dtoUser.AvatarURL = avatar.Url

	return dtoUser, nil
}

func (s *userService) UpdateProfile(ctx context.Context, userID uuid.UUID, req dto.UpdateProfileReq) (dto.User, error) {
	resp, err := s.client.User.UpdateProfile(
		ctx,
		mapper.ToUpdateProfileRequest(userID.String(), req),
	)
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "UpdateProfile", "failed to update user profile", err)
	}

	usr, err := mapper.ToUser(resp)
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "UpdateProfile", "failed to parse to dto.User", err)
	}

	return usr, nil
}

func (s *userService) Follow(ctx context.Context, followerId, followingId uuid.UUID) (err error) {
	_, err = s.client.User.Follow(ctx, &user.FollowRequest{
		FollowerId:  followerId.String(),
		FollowingId: followingId.String(),
	})
	if err != nil {
		return apperror.Wrap("service", "Follow", "failed to follow", err)
	}

	return
}

func (s *userService) Unfollow(ctx context.Context, followerId, followingId uuid.UUID) (err error) {
	_, err = s.client.User.Unfollow(ctx, &user.UnfollowRequest{
		FollowerId:  followerId.String(),
		FollowingId: followingId.String(),
	})

	if err != nil {
		return apperror.Wrap("service", "Unfollow", "failed to unfollow", err)
	}

	return
}

func (s *userService) GetFollowers(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int32,
) (dto.UsersResp, error) {

	resp, err := s.client.User.GetFollowers(
		ctx,
		&user.GetFollowersRequest{
			UserId: userID.String(),
			Page:   page,
			Limit:  limit,
		},
	)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowers", "failed to get followers",
			err,
		)
	}

	usersResp, err := mapper.ToUsersResp(resp, page, limit)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowers", "failed to parse dto.UsersResp",
			err,
		)
	}

	// Collect unique avatar IDs
	avatarSet := make(map[string]struct{})

	for _, u := range usersResp.Users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}

		avatarSet[u.AvatarMediaID.String()] = struct{}{}
	}

	avatarIDs := make([]string, 0, len(avatarSet))

	for id := range avatarSet {
		avatarIDs = append(avatarIDs, id)
	}

	// If nobody has an avatar
	if len(avatarIDs) == 0 {
		return usersResp, nil
	}

	// Get all avatars in one RPC
	mediaResp, err := s.client.Media.GetMedias(
		ctx,
		&media.GetMediasRequest{
			MediaIds: avatarIDs,
		},
	)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowers", "failed to get medias",
			err,
		)
	}

	medias, err := mapper.ToMediaMap(mediaResp)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowers", "failed to map medias",
			err,
		)
	}

	// Attach avatar URL to users
	for i, u := range usersResp.Users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}

		if m, ok := medias[u.AvatarMediaID.String()]; ok {
			usersResp.Users[i].AvatarURL = m.Url
		}
	}

	return usersResp, nil
}

func (s *userService) GetFollowing(
	ctx context.Context,
	userID uuid.UUID,
	page, limit int32,
) (dto.UsersResp, error) {

	resp, err := s.client.User.GetFollowing(ctx, &user.GetFollowingRequest{
		UserId: userID.String(),
		Page:   page,
		Limit:  limit,
	})
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowing", "failed to get following", 
			err,
		)
	}

	usersResp, err := mapper.ToUsersResp(resp, page, limit)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowing", "failed to parse dto.UsersResp",
			err,
		)
	}

	avatarSet := make(map[string]struct{})

	for _, u := range usersResp.Users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}

		avatarSet[u.AvatarMediaID.String()] = struct{}{}
	}

	avatarIDs := make([]string, 0, len(avatarSet))

	for id := range avatarSet {
		avatarIDs = append(avatarIDs, id)
	}

	if len(avatarIDs) == 0 {
		return usersResp, nil
	}

	mediaResp, err := s.client.Media.GetMedias(
		ctx,
		&media.GetMediasRequest{
			MediaIds: avatarIDs,
		},
	)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowing", "failed to get medias",
			err,
		)
	}

	medias, err := mapper.ToMediaMap(mediaResp)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "GetFollowing", "failed to map medias",
			err,
		)
	}

	for i, u := range usersResp.Users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}

		if m, ok := medias[u.AvatarMediaID.String()]; ok {
			usersResp.Users[i].AvatarURL = m.Url
		}
	}

	return usersResp, nil
}

func (s *userService) SearchUsers(
	ctx context.Context,
	req dto.SearchUsersQuery,
) (dto.UsersResp, error) {

	resp, err := s.client.User.SearchUsers(
		ctx,
		&user.SearchUsersRequest{
			Query: req.Query,
			Page:  req.Page,
			Limit: req.Limit,
		},
	)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "SearchUsers", "failed to search users",
			err,
		)
	}

	usersResp, err := mapper.ToUsersResp(resp, req.Page, req.Limit)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "SearchUsers", "failed to parse dto.UsersResp",
			err,
		)
	}

	// Collect unique avatar IDs
	avatarSet := make(map[string]struct{})

	for _, u := range usersResp.Users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}

		avatarSet[u.AvatarMediaID.String()] = struct{}{}
	}

	avatarIDs := make([]string, 0, len(avatarSet))

	for id := range avatarSet {
		avatarIDs = append(avatarIDs, id)
	}

	// No avatars
	if len(avatarIDs) == 0 {
		return usersResp, nil
	}

	// Get all avatars in one RPC
	mediaResp, err := s.client.Media.GetMedias(
		ctx,
		&media.GetMediasRequest{
			MediaIds: avatarIDs,
		},
	)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "SearchUsers", "failed to get medias",
			err,
		)
	}

	medias, err := mapper.ToMediaMap(mediaResp)
	if err != nil {
		return dto.UsersResp{}, apperror.Wrap(
			"service", "SearchUsers", "failed to map medias",
			err,
		)
	}

	// Attach avatar URL
	for i, u := range usersResp.Users {
		if u.AvatarMediaID == uuid.Nil {
			continue
		}

		if m, ok := medias[u.AvatarMediaID.String()]; ok {
			usersResp.Users[i].AvatarURL = m.Url
		}
	}

	return usersResp, nil
}

func (s *userService) GetUserByID(ctx context.Context, userID uuid.UUID) (dto.User, error) {
	resp, err := s.client.User.GetUserById(ctx, &user.GetUserByIdRequest{
		Id: userID.String(),
	})
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "GetUserByID", "failed to get user by id", err)
	}

	usr, err := mapper.ToUser(resp)
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "GetUserByID", "failed to parse dto.User", err)
	}

	avatar, err := s.client.Media.GetMedia(ctx, &media.GetMediaRequest{MediaId: resp.AvatarMediaId})
	if err != nil {
		return dto.User{}, apperror.Wrap("service", "GetUserByID", "failed to get user avatar", err)
	}

	usr.AvatarURL = avatar.Url

	return usr, nil
}
