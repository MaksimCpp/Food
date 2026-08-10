package usecase

import (
	"context"

	"github.com/MaksimCpp/auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	// AccessToken  string
	// RefreshToken string
	// Временно
	Username string
}

type LoginUseCase interface {
	Execute(ctx context.Context, in *LoginInput) (*LoginOutput, error)
}

type PostgreSQLLoginUseCase struct {
	repo domain.UserRepository
}

func NewPostgreSQLLoginUseCase(repo domain.UserRepository) *PostgreSQLLoginUseCase {
	return &PostgreSQLLoginUseCase{
		repo: repo,
	}
}

func (uc *PostgreSQLLoginUseCase) Execute(
	ctx context.Context, in *LoginInput,
) (*LoginOutput, error) {
	user, err := uc.repo.GetByEmail(ctx, in.Email)

	if err != nil {
		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	out := LoginOutput{
		Username: user.Username,
	}

	return &out, nil
}
