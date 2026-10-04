package domain

type TokenService interface {
	GenerateAccessToken(userID int64, role string) (string, error)
	GenerateRefreshToken(userID int64, role string) (string, error)
	ValidateAccessToken(token string) (int64, string, error)
	ValidateRefreshToken(token string) (int64, string, error)
}