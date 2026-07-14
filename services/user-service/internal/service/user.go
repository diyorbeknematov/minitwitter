package service

import (
	"context"
	"log/slog"
	"time"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/diyorbeknematov/minitwitter/gen/go/user"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type userService struct {
	repo   *repository.Repository
	logger *slog.Logger

	user.UnimplementedUserServiceServer
}

func NewUserService(repo *repository.Repository, logger *slog.Logger) *userService {
	return &userService{
		repo:   repo,
		logger: logger,
	}
}

func (s *userService) GetUserById(ctx context.Context, req *user.GetUserByIdRequest) (*user.User, error) {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, apperror.Wrap("service", "GetUserByID", "failed to parse user id", err)
	}

	u, err := s.repo.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap("service", "GetUserById", "failed to get user by id", err)
	}

	return &user.User{
		Id:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Bio:       u.Bio,
		AvatarUrl: u.AvatarMediaID.String(),
		CreatedAt: timestamppb.New(u.CreatedAt),
	}, nil
}

func (s *userService) GetProfile(ctx context.Context, req *user.GetProfileRequest) (*user.User, error) {
	usr, err := s.repo.UserRepo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, apperror.Wrap("service", "GetProfile", "failed to get user profile", err)
	}

	followersCount, err := s.repo.FollowRepo.CountFollowers(ctx, usr.ID)
	if err != nil {
		return nil, apperror.Wrap("service", "GetProfile", "failed to get user followers count", err)
	}

	followingCount, err := s.repo.FollowRepo.CountFollowing(ctx, usr.ID)
	if err != nil {
		return nil, apperror.Wrap("service", "GetProfile", "failed to get user following count", err)
	}

	return &user.User{
		Id:             usr.ID.String(),
		Username:       usr.Username,
		Email:          usr.Email,
		Name:           usr.Name,
		Bio:            usr.Bio,
		AvatarUrl:      usr.AvatarMediaID.String(),
		FollowersCount: uint64(followersCount),
		FollowingCount: uint64(followingCount),
	}, nil
}

func (s *userService) UpdateProfile(ctx context.Context, req *user.UpdateProfileRequest) (*user.User, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, apperror.Wrap("service", "UpdateProfile", "failed to parse user id", err)
	}

	avatarURL, err := uuid.Parse(req.GetAvatarUrl())
	if err != nil {
		return nil, apperror.Wrap("service", "UpdateProfile", "failed to parse avatar media id", err)
	}

	usr, err := s.repo.UserRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap("service", "UpdateProfile", "failed to get user by id", err)
	}

	err = s.repo.UserRepo.Update(ctx, &models.User{
		ID:            usr.ID,
		Username:      usr.Username,
		Email:         usr.Email,
		Name:          req.Name,
		Bio:           req.Bio,
		AvatarMediaID: &avatarURL,
		UpdatedAt:     time.Now(),
	})

	if err != nil {
		return nil, apperror.Wrap("service", "UpdateProfile", "failed to update user profile", err)
	}

	return &user.User{
		Id:        usr.ID.String(),
		Username:  usr.Username,
		Email:     usr.Email,
		Name:      usr.Name,
		Bio:       usr.Bio,
		AvatarUrl: usr.AvatarMediaID.String(),
	}, nil
}

func (s *userService) Follow(ctx context.Context, req *user.FollowRequest) (*user.FollowResponse, error) {
	followerID, err := uuid.Parse(req.FollowerId)
	if err != nil {
		return nil, apperror.Wrap("service", "Follow", "failed to parse follower id", err)
	}

	followingID, err := uuid.Parse(req.FollowingId)
	if err != nil {
		return nil, apperror.Wrap("service", "Follow", "failed to parse following id", err)
	}

	err = s.repo.FollowRepo.Create(ctx, &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, apperror.Wrap("service", "Follow", "failed to follow", err)
	}

	return &user.FollowResponse{
		Success: true,
	}, nil
}

func (s *userService) Unfollow(ctx context.Context, req *user.UnfollowRequest) (*user.UnfollowResponse, error) {
	followerID, err := uuid.Parse(req.FollowerId)
	if err != nil {
		return nil, apperror.Wrap("service", "UnFollow", "failed to parse follower id", err)
	}

	followingID, err := uuid.Parse(req.FollowingId)
	if err != nil {
		return nil, apperror.Wrap("service", "UnFollow", "failed to parse following id", err)
	}

	err = s.repo.FollowRepo.Delete(ctx, &models.Follow{
		FollowerID:  followerID,
		FollowingID: followingID,
		CreatedAt:   time.Now(),
	})
	if err != nil {
		return nil, apperror.Wrap("service", "UnFollow", "failed to follow", err)
	}

	return &user.UnfollowResponse{
		Success: true,
	}, nil
}

func (s *userService) GetFollowers(ctx context.Context, req *user.GetFollowersRequest) (*user.UsersResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, apperror.Wrap("service", "GetFollowers", "failed to parse user id", err)
	}

	users, count, err := s.repo.UserRepo.GetUserFollowers(ctx, userID, int(req.Limit), int(req.Page-1)*int(req.Limit))
	if err != nil {
		return nil, apperror.Wrap("service", "GetFollowers", "failed to get user followers", err)
	}

	return &user.UsersResponse{
		Users: s.toProtoUsers(users),
		Total: uint64(count),
	}, nil
}

func (s *userService) GetFollowing(ctx context.Context, req *user.GetFollowingRequest) (*user.UsersResponse, error) {
	userID, err := uuid.Parse(req.GetUserId())
	if err != nil {
		return nil, apperror.Wrap("service", "GetFollowing", "failed to parse user id", err)
	}

	users, count, err := s.repo.UserRepo.GetUserFollowing(ctx, userID, int(req.Limit), int(req.Page-1)*int(req.Limit))
	if err != nil {
		return nil, apperror.Wrap("service", "GetFollowing", "failed to get user following", err)
	}

	return &user.UsersResponse{
		Users: s.toProtoUsers(users),
		Total: uint64(count),
	}, nil
}

func (s *userService) SearchUsers(ctx context.Context, req *user.SearchUsersRequest) (*user.UsersResponse, error) {
	usrs, total, err := s.repo.UserRepo.Search(ctx, req.String(), int(req.Limit), (int(req.Page)-1)*int(req.Limit))
	if err != nil {
		return nil, apperror.Wrap("service", "SearchUsers", "failed to get users by research", err)
	}

	return &user.UsersResponse{
		Users: s.toProtoUsers(usrs),
		Total: uint64(total),
	}, nil
}

func (s *userService) GetFollowingIds(ctx context.Context, req *user.GetFollowingIdsRequest) (*user.FollowingIdsResponse, error) {
	userID, err := uuid.Parse(req.GetId())
	if err != nil {
		return nil, apperror.Wrap("service", "GetFollowingIds", "failed to parse id", err)
	}

	ids, err := s.repo.FollowRepo.GetFollowingIDs(ctx, userID)
	if err != nil {
		return nil, apperror.Wrap("service", "GetFolliwingIds", "failed to get ids", err)
	}

	respIds := make([]string, 0, len(ids))
	for _, id := range ids {
		respIds = append(respIds, id.String())
	}

	return &user.FollowingIdsResponse{
		Ids: respIds,
	}, nil
}

func (s *userService) toProtoUser(u models.User) *user.User {
	return &user.User{
		Id:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		Name:      u.Name,
		Bio:       u.Bio,
		AvatarUrl: u.AvatarMediaID.String(),
	}
}

func (s *userService) toProtoUsers(users []models.User) []*user.User {
	result := make([]*user.User, 0, len(users))

	for _, u := range users {
		result = append(result, s.toProtoUser(u))
	}

	return result
}
