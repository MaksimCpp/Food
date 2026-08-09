package repository

import (
	"context"
	"errors"

	"github.com/MaksimCpp/auth/internal/domain"
	"github.com/jackc/pgx/v5"
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
		(username, email, password_hash)
		VALUES ($1, $2, $3)
		RETURNING id, email;
	`

	var result domain.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		user.Username,
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

func (repo *PostgreSQLUserRepository) GetByEmail(
	ctx context.Context, email string,
) (*domain.User, error) {
	query := `
		SELECT id, username, email
		FROM users
		WHERE email = $1;
	`

	var result domain.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &result, nil
}