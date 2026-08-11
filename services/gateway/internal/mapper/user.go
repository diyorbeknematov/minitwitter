package mapper

import (
	"github.com/diyorbeknematov/minitwitter/gen/go/user"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/google/uuid"
)

func ToUser(pb *user.User) (dto.User, error) {
	userId, err := uuid.Parse(pb.Id)
	if err != nil {
		return dto.User{}, err
	}

	return dto.User{
		ID:             userId,
		Username:       pb.Username,
		Email:          pb.Email,
		Name:           pb.Name,
		Bio:            pb.Bio,
		AvatarURL:      pb.AvatarUrl,
		FollowersCount: pb.FollowersCount,
		FollowingCount: pb.FollowingCount,
		CreatedAt:      pb.CreatedAt.AsTime(),
	}, nil
}

func ToUsersResp(
	pb *user.UsersResponse,
	page, limit int32,
) (dto.UsersResp, error) {

	users := make([]dto.User, len(pb.Users))

	for i, u := range pb.Users {
		user, err := ToUser(u)
		if err != nil {
			return dto.UsersResp{}, err
		}
		users[i] = user
	}

	return dto.UsersResp{
		Users: users,
		Pagination: dto.Pagination{
			Page:  page,
			Limit: limit,
			Total: int64(pb.Total),
		},
	}, nil
}

func ToUpdateProfileRequest(
	userID string,
	req dto.UpdateProfileReq,
) *user.UpdateProfileRequest {

	return &user.UpdateProfileRequest{
		UserId:    userID,
		Name:      req.Name,
		Bio:       req.Bio,
		AvatarUrl: req.AvatarURL,
	}
}
