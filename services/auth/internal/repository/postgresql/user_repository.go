package repository

import (
	"context"
	"errors"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLUserRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLUserRepository(pool *pgxpool.Pool) *PostgreSQLUserRepository {
	return &PostgreSQLUserRepository{
		pool: pool,
	}
}

	// Create(ctx context.Context, user *User) (*User, error)
	// GetByEmail(ctx context.Context, email string) (*User, error)
	// GetByID(ctx context.Context, id int64) (*User, error)

func (repo *PostgreSQLUserRepository) Create(
	ctx context.Context, user *domain.User,
) (*domain.User, error) {
	query := `
		INSERT INTO users
		(email, password_hash)
		VALUES ($1, $2)
		RETURNING id, email;
	`

	var result domain.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		user.Email,
		user.PasswordHash,
	).Scan(
		&result.ID,
		&result.Email,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return nil, domain.ErrUserAlreadyExist
		}

		return nil, err
	}

	return &result, nil
}