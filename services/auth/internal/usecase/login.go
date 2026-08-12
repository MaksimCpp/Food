package usecase

import (
	"context"
	"errors"

	"github.com/MaksimCpp/auth/internal/domain"
	"golang.org/x/crypto/bcrypt"
)

type LoginInput struct {
	Email    string
	Password string
}

type LoginOutput struct {
	AccessToken  string
	RefreshToken string
}

type LoginUseCase interface {
	Execute(ctx context.Context, in *LoginInput) (*LoginOutput, error)
}

type PostgreSQLLoginUseCase struct {
	repo         domain.UserRepository
	tokenService domain.TokenService
}

func NewPostgreSQLLoginUseCase(
	repo domain.UserRepository,
	tokenService domain.TokenService,
) *PostgreSQLLoginUseCase {
	return &PostgreSQLLoginUseCase{
		repo: repo,
		tokenService: tokenService,
	}
}

func (uc *PostgreSQLLoginUseCase) Execute(
	ctx context.Context, in *LoginInput,
) (*LoginOutput, error) {
	user, err := uc.repo.GetByEmail(ctx, in.Email)

	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return nil, domain.ErrInvalidCredentials
		}

		return nil, err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(in.Password))

	if err != nil {
		return nil, domain.ErrInvalidCredentials
	}

	accessToken, err := uc.tokenService.GenerateAccessToken(user.ID)

	if err != nil {
		return nil, err
	}

	refreshToken, err := uc.tokenService.GenerateRefreshToken(user.ID)

	if err != nil {
		return nil, err
	}

	out := LoginOutput{
		AccessToken: accessToken,
		RefreshToken: refreshToken,
	}

	return &out, nil
}
