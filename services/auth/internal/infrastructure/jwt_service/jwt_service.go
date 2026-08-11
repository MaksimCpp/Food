package jwtservice

import (
	"time"

	"github.com/MaksimCpp/auth/internal/config"
	"github.com/golang-jwt/jwt/v5"
)

type JWTTokenService struct {
	secretKey string
	ttl       time.Duration
}

type Payload struct {
	jwt.RegisteredClaims
	UserID int64
}

	// GenerateAccessToken(userID int64) (string, error)
	// GenerateRefreshToken(userID int64) (string, error)

func NewJWTTokenService(secretKey string, ttl time.Duration) *JWTTokenService {
	return &JWTTokenService{
		secretKey: secretKey,
		ttl: ttl,
	}
}

func (s *JWTTokenService) GenerateAccessToken(userID int64) (string, error) {
	payload := Payload{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(s.ttl)),
		},
	}

	token := jwt.NewWithClaims(
		jwt.SigningMethodHS256, 
		&payload,
	)

	return token.SignedString([]byte(s.secretKey))
}