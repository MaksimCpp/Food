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

func (repo *PostgreSQLUserRepository) Create(
	ctx context.Context, user *domain.User,
) (*domain.User, error) {
	query := `
		INSERT INTO users
		(username, email, role, password_hash)
		VALUES ($1, $2, $3, $4)
		RETURNING id, username, email, role, password_hash, created_at;
	`

	var result domain.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		user.Username,
		user.Email,
		user.Role,
		user.PasswordHash,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.Role,
		&result.PasswordHash,
		&result.CreatedAt,
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
		SELECT id, username, email, role, password_hash, created_at
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
		&result.Role,
		&result.PasswordHash,
		&result.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLUserRepository) GetByID(
	ctx context.Context, id int64,
) (*domain.User, error) {
	query := `
		SELECT id, username, email, role, password_hash, created_at
		FROM users
		WHERE id = $1;
	`

	var result domain.User

	err := repo.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&result.ID,
		&result.Username,
		&result.Email,
		&result.Role,
		&result.PasswordHash,
		&result.CreatedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrUserNotFound
		}

		return nil, err
	}

	return &result, nil
}