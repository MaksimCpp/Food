package usecase

import (
	"context"
	"errors"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgconn"
	"golang.org/x/crypto/bcrypt"
)

type RegisterInput struct {
	Email    string
	Password string
}

type RegisterOutput struct {
	ID int64
	Email  string
}

type RegisterUseCase interface {
	Execute(ctx context.Context, in *RegisterInput) (*RegisterOutput, error)
}

type PostgreSQLRegisterUseCase struct {
	repo domain.UserRepository
}

func NewPostgreSQLRegisterUseCase(repo domain.UserRepository) *PostgreSQLRegisterUseCase {
	return &PostgreSQLRegisterUseCase{
		repo: repo,
	}
}

func (uc *PostgreSQLRegisterUseCase) Execute(
	ctx context.Context, in *RegisterInput,
) (*RegisterOutput, error) {
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(in.Password), 10)

	if err != nil {
		return nil, err
	}

	user := domain.User{
		Email: in.Email,
		PasswordHash: string(passwordHash),
	}

	result, err := uc.repo.Create(ctx, &user)

	if err != nil {
		var pgErr *pgconn.PgError
	
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrUserAlreadyExist
		}
		
		return nil, err
	}

	out := RegisterOutput{
		ID: result.ID,
		Email: result.Email,
	}

	return &out, nil
}
