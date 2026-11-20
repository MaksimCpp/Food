package repository

import (
	"context"
	"errors"

	"github.com/MaksimCpp/orders/internal/domain"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgreSQLOrdersRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLOrdersRepository(pool *pgxpool.Pool) *PostgreSQLOrdersRepository {
	return &PostgreSQLOrdersRepository{
		pool: pool,
	}
}

// CREATE TYPE order_status AS ENUM (
//     'accepted',
//     'out_for_delivery',
//     'delivered'
// );

//     user_id INTEGER NOT NULL CHECK (user_id > 0),
//     product_id INTEGER NOT NULL CHECK (product_id > 0),
//     products_quantity INTEGER NOT NULL CHECK (products_quantity > 0),
//     courier_id INTEGER CHECK (courier_id > 0),
//     address VARCHAR(400) NOT NULL, -- в будущем добавить GEOGRAPHY
//     status order_status NOT NULL DEFAULT 'ACCEPT',

func (repo *PostgreSQLOrdersRepository) CreateOrder(
	ctx context.Context, order *domain.Order,
) (*domain.Order, error) {
	query := `
		INSERT INTO 
			orders (
				user_id, 
				product_id,
				products_quantity,
				address,
				status
			)
		VALUES 
			($1, $2, $3, $4, 'accepted')
		RETURNING
			id, user_id, product_id, products_quantity, address, status;
	`

	result := domain.Order{}

	err := repo.pool.QueryRow(
		ctx,
		query,
		order.UserID,
		order.ProductID,
		order.ProductsQuantity,
		order.Address,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.ProductID,
		&result.ProductsQuantity,
		&result.Address,
		&result.Status,
	)

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo * PostgreSQLOrdersRepository) ListOrdersWithoutCourier(
	ctx context.Context,
) ([]*domain.Order, error) {
	query := `
		SELECT
			id,
			user_id, 
			product_id,
			products_quantity,
			COALESCE(courier_id, 0) AS courier_id,
			address,
			status
		FROM
			orders
		WHERE 
			courier_id IS NULL;
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := []*domain.Order{}

	for rows.Next() {
		order := domain.Order{}

		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.ProductID,
			&order.ProductsQuantity,
			&order.CourierID,
			&order.Address,
			&order.Status,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PostgreSQLOrdersRepository) ListOrdersByCourierID(
	ctx context.Context, courierID int64,
) ([]*domain.Order, error) {
	query := `
		SELECT
			id,
			user_id, 
			product_id,
			products_quantity,
			COALESCE(courier_id, 0) AS courier_id,
			address,
			status
		FROM
			orders
		WHERE 
			courier_id = $1;
	`

	rows, err := repo.pool.Query(
		ctx,
		query,
		courierID,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	results := []*domain.Order{}

	for rows.Next() {
		order := domain.Order{}

		err = rows.Scan(
			&order.ID,
			&order.UserID,
			&order.ProductID,
			&order.ProductsQuantity,
			&order.CourierID,
			&order.Address,
			&order.Status,
		)

		if err != nil {
			return nil, err
		}

		results = append(results, &order)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return results, nil
}

func (repo *PostgreSQLOrdersRepository) GetOrder(
	ctx context.Context, id int64,
) (*domain.Order, error) {
	query := `
		SELECT
			id,
			user_id, 
			product_id,
			products_quantity,
			COALESCE(courier_id, 0) AS courier_id,
			address,
			status
		FROM
			orders
		WHERE 
			id = $1;
	`

	var result domain.Order

	err := repo.pool.QueryRow(
		ctx,
		query,
		id,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.ProductID,
		&result.ProductsQuantity,
		&result.CourierID,
		&result.Address,
		&result.Status,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, domain.ErrOrderNotFound
		}
		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLOrdersRepository) AcceptOrder(
	ctx context.Context, id int64, courierID int64,
) (*domain.Order, error) {
	query := `
		UPDATE
			orders
		SET
			courier_id = $1,
			status = $2
		WHERE 
			id = $3 AND courier_id IS NULL AND status = $4
		RETURNING
			id,
			user_id, 
			product_id,
			products_quantity,
			COALESCE(courier_id, 0) AS courier_id,
			address,
			status;
	`

	var result domain.Order

	err := repo.pool.QueryRow(
		ctx,
		query,
		courierID,
		domain.StatusOutForDelivery,
		id,
		domain.StatusAccepted,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.ProductID,
		&result.ProductsQuantity,
		&result.CourierID,
		&result.Address,
		&result.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		order, errNotFound := repo.GetOrder(ctx, id)

		if errNotFound != nil {
			return nil, errNotFound
		}

		if order.CourierID != 0 {
			return nil, domain.ErrOrderAlreadyAccept
		}

		return nil, domain.ErrInvalidOrderTransition
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (repo *PostgreSQLOrdersRepository) CompleteOrder(
	ctx context.Context, id int64,
) (*domain.Order, error) {
	query := `
		UPDATE
			orders
		SET
			status = $1
		WHERE 
			id = $2 AND courier_id IS NOT NULL AND status = $3
		RETURNING
			id,
			user_id, 
			product_id,
			products_quantity,
			COALESCE(courier_id, 0) AS courier_id,
			address,
			status;
	`

	var result domain.Order

	err := repo.pool.QueryRow(
		ctx,
		query,
		domain.StatusDelivered,
		id,
		domain.StatusOutForDelivery,
	).Scan(
		&result.ID,
		&result.UserID,
		&result.ProductID,
		&result.ProductsQuantity,
		&result.CourierID,
		&result.Address,
		&result.Status,
	)

	if errors.Is(err, pgx.ErrNoRows) {
		_, errNotFound := repo.GetOrder(ctx, id)

		if errNotFound != nil {
			return nil, errNotFound
		}

		return nil, domain.ErrInvalidOrderTransition
	}

	if err != nil {
		return nil, err
	}

	return &result, nil
}