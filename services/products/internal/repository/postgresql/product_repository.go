package repository

import (
	"context"
	"errors"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/MaksimCpp/products/internal/dto"
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
		INSERT INTO products (name, description, price, image_key, is_available, stock_quantity)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, name, description, price, image_key, is_available, stock_quantity;
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
		product.StockQuantity,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.Price,
		&result.ImageKey,
		&result.IsAvailable,
		&result.StockQuantity,
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
		SELECT id, name, description, price, image_key, is_available, stock_quantity
		FROM products
		WHERE id = $1;
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
		&result.StockQuantity,
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
		SELECT id, name, description, price, image_key, is_available, stock_quantity
		FROM products
		LIMIT $1
		OFFSET $2;
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
		result := &domain.Product{}

		err := rows.Scan(
			&result.ID,
			&result.Name,
			&result.Description,
			&result.Price,
			&result.ImageKey,
			&result.IsAvailable,
			&result.StockQuantity,
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
			name = $1,
			description = $2,
			price = $3,
			image_key = $4,
			is_available = $5,
			stock_quantity = $6
		WHERE
			id = $7
		RETURNING
			id, name, description, price, image_key, is_available, stock_quantity;
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
		newProduct.StockQuantity,
		id,
	).Scan(
		&result.ID,
		&result.Name,
		&result.Description,
		&result.Price,
		&result.ImageKey,
		&result.IsAvailable,
		&result.StockQuantity,
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
		WHERE id = $1;
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
		LIMIT $1
		OFFSET $2;
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
	ctx context.Context, ids ...int64,
) ([]*domain.Product, error) {
	query := `
		SELECT id, name, description, price, image_key, is_available, stock_quantity
		FROM products
		WHERE id = ANY($1);
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
			&product.StockQuantity,
		)

		if err != nil {
			return nil, err
		}

		result = append(result, &product)
	}

	return result, nil
}

func (repo *PostgreSQLProductsRepository) Reserve(
	ctx context.Context, in *dto.ReserveProductInputOutput,
) (*dto.ReserveProductInputOutput, error) {
	query := `
		UPDATE products
		SET reserved_quantity = reserved_quantity + $1
		WHERE id = $2 AND reserved_quantity + $1 <= stock_quantity
		RETURNING id, $1;
	`

	result := dto.ReserveProductInputOutput{}

	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		query,
		in.Quantity,
		in.ID,
	).Scan(
		&result.ID,
		&result.Quantity,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var exist bool
			err = tx.QueryRow(
				ctx,
				`SELECT EXISTS(SELECT 1 FROM products WHERE id = $1);`,
				in.ID,
			).Scan(&exist)

			if err != nil {
				return nil, err
			} 
			
			if !exist {
				return nil, domain.ErrProductNotFound
			}

			return nil, domain.ErrProductQuantity
		}
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLProductsRepository) DeductDeliveredStock(
	ctx context.Context, in *dto.DeductDeliveredStockInput,
) (*dto.DeductDeliveredStockOutput, error) {
	query := `
		UPDATE products
		SET stock_quantity = stock_quantity - $1,
			reserved_quantity = reserved_quantity - $1
		WHERE id = $2 AND stock_quantity >= $1 AND reserved_quantity >= $1
		RETURNING id, stock_quantity;
	`

	result := dto.DeductDeliveredStockOutput{}

	tx, err := repo.pool.BeginTx(ctx, pgx.TxOptions{
		IsoLevel: pgx.ReadCommitted,
	})

	if err != nil {
		return nil, err
	}

	defer tx.Rollback(ctx)

	err = tx.QueryRow(
		ctx,
		query,
		in.Quantity,
		in.ID,
	).Scan(
		&result.ID,
		&result.StockQuantity,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			var exist bool
			err = tx.QueryRow(
				ctx,
				`SELECT EXISTS(SELECT 1 FROM products WHERE id = $1);`,
				in.ID,
			).Scan(&exist)

			if err != nil {
				return nil, err
			} 
			
			if !exist {
				return nil, domain.ErrProductNotFound
			}

			return nil, domain.ErrProductQuantity
		}
		return nil, err
	}

	err = tx.Commit(ctx)

	if err != nil {
		return nil, err
	}

	return &result, nil
}