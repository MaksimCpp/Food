package domain

type TokenService interface {
	GenerateAccessToken(userID int64) (string, error)
	GenerateRefreshToken(userID int64) (string, error)
	ValidateAccessToken(token string) (int64, error)
}