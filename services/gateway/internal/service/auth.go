package service

import (
	"context"

	"github.com/diyorbeknematov/minitwitter/gen/go/auth"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/grpcclient"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/mapper"
	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
)

type authService struct {
	client *grpcclient.Client
}

func NewAuthService(c *grpcclient.Client) *authService {
	return &authService{
		client: c,
	}
}

func (s *authService) Register(ctx context.Context, req dto.RegisterReq) (dto.RegisterResp, error) {
	resp, err := s.client.Auth.Register(
		ctx,
		mapper.ToRegisterReq(req),
	)
	if err != nil {
		return dto.RegisterResp{}, apperror.Wrap("service", "Register", "failed to register user", err)
	}

	return dto.RegisterResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
	}, nil
}

func (s *authService) Login(ctx context.Context, req dto.LoginReq) (dto.LoginResp, error) {
	resp, err := s.client.Auth.Login(
		ctx,
		mapper.ToLoginReq(req),
	)
	if err != nil {
		return dto.LoginResp{}, apperror.Wrap("service", "Login", "failed to login user", err)
	}

	return dto.LoginResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt.AsTime(),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req dto.RefreshTokenReq) (dto.LoginResp, error) {
	resp, err := s.client.Auth.RefreshToken(ctx, &auth.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	})
	if err != nil {
		return dto.LoginResp{}, apperror.Wrap("service", "RefreshToken", "failed to refresh token", err)
	}

	return dto.LoginResp{
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresAt:    resp.ExpiresAt.AsTime(),
	}, nil
}
