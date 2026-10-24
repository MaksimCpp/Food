package repository

import (
	"context"
	"errors"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLProductsRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLProductsRepository(pool *pgxpool.Pool) *PostgreSQLProductsRepository {
	return &PostgreSQLProductsRepository{
		pool: pool,
	}
}

func (repo *PostgreSQLProductsRepository) Create(
	ctx context.Context, product *domain.Product,
) (*domain.Product, error) {
	query := `
		INSERT INTO products (name, description, price, image_key, is_available)
		VALUES (%1, %2, %3, %4, %5)
		RETURNING name, description, price, image_key, is_available;
	`

	var result domain.Product

	err := repo.pool.QueryRow(
		ctx,
		query,
		product.Name,
		product.Description,
		product.Price,
		product.ImageKey,
		product.IsAvailable,
	).Scan(
		&result.Name,
		&result.Description,
		&result.Price,
		&result.ImageKey,
		&result.IsAvailable,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLProductsRepository) Get(
	ctx context.Context, id int64,
) (*domain.Product, error) {
	query := `
		SELECT id, name, description, price, image_key, is_available
		FROM products
		WHERE id = %1;
	`

	var result domain.Product

	err := repo.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.Price,
		&result.ImageKey,
		&result.IsAvailable,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLProductsRepository) List(
	ctx context.Context, limit int, offset int,
) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, price, image_key, is_available
		FROM products
		LIMIT %1
		OFFSET %2;
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var results []*domain.Product

	for rows.Next() {
		var result *domain.Product

		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.Description,
			&result.Price,
			&result.ImageKey,
			&result.IsAvailable,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, result)
	}

	return results, nil
}

func (repo *PostgreSQLProductsRepository) Update(
	ctx context.Context, newProduct *domain.Product, id int64,
) (*domain.Product, error) {
	query := `
		UPDATE 
			products
		SET
			name = %1,
			description = %2,
			price = %3,
			image_key = %4,
			is_available = %5
		WHERE
			id = %6
		RETURNING
			id, name, description, price, image_key, is_available;
	`

	var result domain.Product

	err := repo.pool.QueryRow(
		ctx,
		query,
		newProduct.Name,
		newProduct.Description,
		newProduct.Price,
		newProduct.ImageKey,
		newProduct.IsAvailable,
		id,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.Price,
		&result.ImageKey,
		&result.IsAvailable,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrProductNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLProductsRepository) Delete(
	ctx context.Context, id int64,
) error {
	query := `
		DELETE FROM products
		WHERE id = %1;
	`

	ct, err := repo.pool.Exec(
		ctx,
		query,
		id,
	)

	if err != nil {
		return err
	}

	if ct.RowsAffected() == 0 {
		return domain.ErrProductNotFound
	}

	return nil
}

func (repo *PostgreSQLProductsRepository) GetIDs(
	ctx context.Context, limit int, offset int,
) ([]*domain.Product, error) {
	query := `
		SELECT id
		FROM products
		LIMIT %1
		OFFSET %2;
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*domain.Product

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(&product.ID)

		if err != nil {
			return nil, err
		}

		result = append(result, &product)
	}

	return result, nil
}

func (repo *PostgreSQLProductsRepository) GetByIDs(
	ctx context.Context, ids... int64,
) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, price, image_key, is_available
		FROM products
		WHERE id = ANY(%1);
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
		ids,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var result []*domain.Product

	for rows.Next() {
		var product domain.Product
		err := rows.Scan(
			&product.ID,
			&product.Name,
			&product.Description,
			&product.Price,
			&product.ImageKey,
			&product.IsAvailable,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, &product)
	}

	return result, nil
}