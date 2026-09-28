package jwtservice

import (
	"time"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenTypeAccess  = "access"
	TokenTypeRefresh = "refresh"
)

type JWTTokenService struct {
	secretKey   string
	access_ttl  time.Duration
	refresh_ttl time.Duration
}

type Payload struct {
	jwt.RegisteredClaims
	UserID    int64
	TokenType string
}

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
		TokenType: TokenTypeAccess,
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
		TokenType: TokenTypeRefresh,
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

func (s *JWTTokenService) ValidateAccessToken(token string) (int64, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token, 
		&Payload{}, 
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, domain.ErrInvalidToken
			}

			return s.secretKey, nil
		},
	)

	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	if !jwtToken.Valid {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return 0, domain.ErrInvalidToken
	}

	if claims.TokenType != TokenTypeAccess {
		return 0, domain.ErrInvalidToken
	}

	return claims.UserID, nil
}

func (s *JWTTokenService) ValidateRefreshToken(token string) (int64, error) {
	jwtToken, err := jwt.ParseWithClaims(
		token, 
		&Payload{}, 
		func(t *jwt.Token) (interface{}, error) {
			if t.Method != jwt.SigningMethodHS256 {
				return nil, domain.ErrInvalidToken
			}

			return s.secretKey, nil
		},
	)

	if err != nil {
		return 0, domain.ErrInvalidToken
	}

	if !jwtToken.Valid {
		return 0, domain.ErrInvalidToken
	}

	claims, ok := jwtToken.Claims.(*Payload)

	if !ok {
		return 0, domain.ErrInvalidToken
	}

	if claims.TokenType != TokenTypeRefresh {
		return 0, domain.ErrInvalidToken
	}

	return claims.UserID, nil
}