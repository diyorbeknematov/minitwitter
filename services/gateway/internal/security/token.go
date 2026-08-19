package security

import (
	"errors"

	"github.com/diyorbeknematov/minitwitter/services/gateway/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
)

type jwtCustomClaim struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func ValidateToken(tokenString, secret string) (*jwtCustomClaim, error) {

	token, err := jwt.ParseWithClaims(
		tokenString,
		&jwtCustomClaim{},
		func(t *jwt.Token) (any, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}

			return []byte(secret), nil
		},
	)

	if err != nil {
		return nil, apperror.Wrap("security", "ValidateToken", "failed to parse token", err)
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	claims, ok := token.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("invalid token")
	}

	return claims, nil
}
