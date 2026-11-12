package repository

import "github.com/jackc/pgx/v5/pgxpool"

type PostgreSQLOrdersRepository struct {
	pool *pgxpool.Pool
}

func NewPostgreSQLOrdersRepository(pool *pgxpool.Pool) *PostgreSQLOrdersRepository {
	return &PostgreSQLOrdersRepository{
		pool: pool,
	}
}


//  Create(ctx context.Context, order *Order) (*Order, error)
// 	List(ctx context.Context) ([]*Order, error)
// 	Get(ctx context.Context, id int64) (*Order, error)
// 	Accept(ctx context.Context, id int64, courier_id int64) (*Order, error)
// 	Complete(ctx context.Context, id int64) (*Order, error)