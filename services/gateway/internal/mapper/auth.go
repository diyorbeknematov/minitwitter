package mapper

import (
	"github.com/diyorbeknematov/minitwitter/gen/go/auth"
	"github.com/diyorbeknematov/minitwitter/services/gateway/internal/dto"
)

func ToRegisterReq(req dto.RegisterReq) *auth.RegisterRequest {
	return &auth.RegisterRequest{
		Username: req.Username,
		Email:    req.Email,
		Name:     req.Name,
		Password: req.Password,
	}
}

func ToLoginReq(req dto.LoginReq) *auth.LoginRequest {
	return &auth.LoginRequest{
		Email:    req.Email,
		Password: req.Password,
	}
}
