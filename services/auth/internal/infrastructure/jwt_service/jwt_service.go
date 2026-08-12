package jwtservice

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenService struct {
	secretKey   string
	access_ttl  time.Duration
	refresh_ttl time.Duration
}

type Payload struct {
	jwt.RegisteredClaims
	UserID int64
}

	// GenerateAccessToken(userID int64) (string, error)
	// GenerateRefreshToken(userID int64) (string, error)

func NewJWTTokenService(
	secretKey string, 
	access_ttl  time.Duration,
	refresh_ttl time.Duration,
) *JWTTokenService {
	return &JWTTokenService{
		secretKey: secretKey,
		access_ttl: access_ttl,
		refresh_ttl: refresh_ttl,
	}
}

func (s *JWTTokenService) GenerateAccessToken(userID int64) (string, error) {
	payload := Payload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.access_ttl)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, 
		&payload,
	)

	return token.SignedString([]byte(s.secretKey))
}

func (s *JWTTokenService) GenerateRefreshToken(userID int64) (string, error) {
	payload := Payload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.refresh_ttl)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, 
		&payload,
	)

	return token.SignedString([]byte(s.secretKey))
}