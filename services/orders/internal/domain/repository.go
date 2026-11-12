package domain

import (
	"context"
)

type OrdersRepository interface {
	Create(ctx context.Context, order *Order) (*Order, error)
	List(ctx context.Context) ([]*Order, error)
	Get(ctx context.Context, id int64) (*Order, error)
	Accept(ctx context.Context, id int64, courier_id int64) (*Order, error)
	Complete(ctx context.Context, id int64) (*Order, error)
}