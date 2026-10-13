package repository

import (
	"context"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLProductRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLProductRepository(pool *pgxpool.Pool) *PostgreSQLProductRepository {
	return &PostgreSQLProductRepository{
		pool: pool,
	}
}

	// Get(ctx context.Context, id int64) (*Product, error)
	// List(ctx context.Context, limit int, offset int) ([]*Product, error)
	// Update(ctx context.Context, newProduct *Product, id int64) (*Product, error)
	// Delete(ctx context.Context, id int64) (*Product, error)

func (repo *PostgreSQLProductRepository) Create(
	ctx context.Context, product *domain.Product,
) (*domain.Product, error) {
	query := `
		INSERT INTO products (name, description, price, image_url, is_available)
		VALUES (%1, %2, %3, %4, %5)
		RETURNING name, description, price, image_url, is_available;
	`

	var result domain.Product

	err := repo.pool.QueryRow(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.ImageUrl,
		product.IsAvailable,
	).Scan(
		&result.Name,
		&result.Description,
		&result.Price,
		&result.ImageUrl,
		&result.IsAvailable,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}