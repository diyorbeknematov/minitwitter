package service

import (
	"context"
	"database/sql"
	"errors"
	"log/slog"

	"github.com/diyorbek/minitwitter/services/user-service/internal/config"
	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/internal/repository"
	"github.com/diyorbek/minitwitter/services/user-service/internal/security"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/diyorbeknematov/minitwitter/gen/go/auth"
	"github.com/google/uuid"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type authService struct {
	repo   *repository.Repository
	cfg    *config.Config
	logger *slog.Logger

	auth.UnimplementedAuthServiceServer
}

func NewAuthService(repo *repository.Repository, cfg *config.Config, logger *slog.Logger) *authService {
	return &authService{
		repo:   repo,
		cfg:    cfg,
		logger: logger,
	}
}

func (s *authService) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
	_, err := s.repo.UserRepo.GetByEmail(ctx, req.GetEmail())
	if err == nil {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "the email is already exists", apperror.ErrEmailExists)
	} else if !errors.Is(err, sql.ErrNoRows) {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "failed to check email exists", err)
	}

	_, err = s.repo.UserRepo.GetByUsername(ctx, req.GetUsername())
	if err == nil {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "the username is already exists", errors.New("username already exists"))
	} else if !errors.Is(err, sql.ErrNoRows) {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "failed to check username exists", err)
	}

	hashedPassword, err := security.HashPassword(req.GetPassword())
	if err != nil {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "failed to hash user", err)
	}
	user := &models.User{
		ID:           uuid.New(),
		Username:     req.GetUsername(),
		Email:        req.GetEmail(),
		PasswordHash: hashedPassword,
		Name:         req.GetName(),
	}
	err = s.repo.UserRepo.Create(ctx, user)

	if err != nil {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "failed to create user", err)
	}

	accessToken, refreshToken, err := s.createTokenPair(user.ID, user.Username)
	if err != nil {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "failed to generate tokens", err)
	}

	err = s.saveRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return &auth.RegisterResponse{}, apperror.Wrap("service", "Register", "failed to save refresh token to database", err)
	}

	return &auth.RegisterResponse{
		AccessToken:  accessToken.TokenStr,
		RefreshToken: refreshToken.TokenStr,
	}, nil
}

func (s *authService) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
	user, err := s.repo.UserRepo.GetByEmail(ctx, req.GetEmail())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &auth.LoginResponse{}, apperror.Wrap("service", "Login", "User not found to get by email", apperror.ErrNotFound)
		}

		return &auth.LoginResponse{}, apperror.Wrap("service", "Login", "failed to check if user exists", err)
	}

	if !security.VerifyPassword(req.GetPassword(), user.PasswordHash) {
		return &auth.LoginResponse{}, apperror.Wrap("service", "Login", "failed to check password", apperror.ErrInvalidInput)
	}

	accessToken, refreshToken, err := s.createTokenPair(user.ID, user.Username)
	if err != nil {
		return &auth.LoginResponse{}, apperror.Wrap("service", "Login", "failed to generate tokens", err)
	}

	err = s.repo.RefreshTokenRepo.DeleteByUserID(ctx, user.ID)
	if err != nil {
		return &auth.LoginResponse{}, apperror.Wrap("service", "Login", "failed to delete refresh token", err)
	}

	err = s.saveRefreshToken(ctx, user.ID, refreshToken)
	if err != nil {
		return &auth.LoginResponse{}, apperror.Wrap("service", "Login", "failed to save refresh token to database", err)
	}

	return &auth.LoginResponse{
		AccessToken:  accessToken.TokenStr,
		RefreshToken: refreshToken.TokenStr,
		ExpiresAt:    timestamppb.New(accessToken.ExpiresAt),
	}, nil
}

func (s *authService) RefreshToken(ctx context.Context, req *auth.RefreshTokenRequest) (*auth.LoginResponse, error) {
	token, err := security.ParseToken(req.RefreshToken, s.cfg.RefreshtokenSecret)
	if err != nil {
		return &auth.LoginResponse{}, apperror.Wrap("service", "RefreshToken", "failed to refresh token", err)
	}

	refreshToken, err := s.repo.RefreshTokenRepo.GetByUserID(ctx, token.UserID)
	if err != nil {
		return &auth.LoginResponse{}, apperror.Wrap("service", "RefreshToken", "failed to get refresh token from database by user id", err)
	}

	if !security.VerifyToken(req.RefreshToken, refreshToken.TokenHash) {
		return &auth.LoginResponse{}, apperror.Wrap("service", "RerfeshToken", "failed to verify refresh token", err)
	}

	accessToken, err := security.GenerateAccessToken(token.UserID, token.Username, s.cfg.AccesstokenSecret)
	if err != nil {
		return &auth.LoginResponse{}, apperror.Wrap("service", "RefreshToken", "failed to generate a new access token", err)
	}

	return &auth.LoginResponse{
		AccessToken:  accessToken.TokenStr,
		RefreshToken: req.RefreshToken,
	}, nil
}

func (s *authService) createTokenPair(userID uuid.UUID, username string) (*models.Token, *models.Token, error) {
	accessToken, err := security.GenerateAccessToken(
		userID,
		username,
		s.cfg.AccesstokenSecret,
	)
	if err != nil {
		return nil, nil, err
	}

	refreshToken, err := security.GenerateRefreshToken(
		userID,
		username,
		s.cfg.RefreshtokenSecret,
	)
	if err != nil {
		return nil, nil, err
	}

	return accessToken, refreshToken, nil
}

func (s *authService) saveRefreshToken(ctx context.Context, userID uuid.UUID, token *models.Token) error {
	hashedToken, err := security.HashToken(token.TokenStr)
	if err != nil {
		return apperror.Wrap("service", "saveRefreshToken", "failed to hash refresh token", err)
	}

	return s.repo.RefreshTokenRepo.Create(ctx, models.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: hashedToken,
		ExpiresAt: token.ExpiresAt,
	})
}
