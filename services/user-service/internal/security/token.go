package security

import (
	"errors"
	"time"

	"github.com/diyorbek/minitwitter/services/user-service/internal/models"
	"github.com/diyorbek/minitwitter/services/user-service/pkg/apperror"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

type jwtCustomClaim struct {
	UserID   uuid.UUID `json:"user_id"`
	Username string    `json:"username"`
	jwt.RegisteredClaims
}

func genToken(userID uuid.UUID, username, secure string, expiresAt time.Time) (*models.Token, error) {
	claims := jwtCustomClaim{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	token, err := jwtToken.SignedString([]byte(secure))
	if err != nil {
		return nil, err
	}

	return &models.Token{
		TokenStr:  token,
		ExpiresAt: expiresAt,
	}, nil
}

func GenerateAccessToken(userID uuid.UUID, username, secure string) (*models.Token, error) {
	accessToken, err := genToken(userID, username, secure, time.Now().Add(3*time.Hour))
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func GenerateRefreshToken(userID uuid.UUID, username, secure string) (*models.Token, error) {
	accessToken, err := genToken(userID, username, secure, time.Now().Add(7*24*time.Hour))
	if err != nil {
		return nil, err
	}

	return accessToken, nil
}

func ParseToken(token, secure string) (*jwtCustomClaim, error) {
	jwtToken, err := jwt.ParseWithClaims(token, &jwtCustomClaim{}, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid signing method")
		}

		return []byte(secure), nil
	})
	if err != nil {
		return nil, apperror.Wrap("service", "ParseToken", "failed to parse token", err)
	}

	claims, ok := jwtToken.Claims.(*jwtCustomClaim)
	if !ok {
		return nil, errors.New("token claims are not of type *jwtCustomClaim")
	}

	return claims, nil
}

func TokenValid(tokenString, secure string) bool {
	_, err := ParseToken(tokenString, secure)
	return err == nil
}
