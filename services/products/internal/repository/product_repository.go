package repository

import "github.com/jackc/pgx/v5/pgxpool"

type PostgreSQLProductRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLProductRepository(pool *pgxpool.Pool) *PostgreSQLProductRepository {
	return &PostgreSQLProductRepository{
		pool: pool,
	}
}

