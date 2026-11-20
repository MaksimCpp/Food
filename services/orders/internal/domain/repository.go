package domain

import (
	"context"
)

type OrdersRepository interface {
	CreateOrder(ctx context.Context, order *Order) (*Order, error)
	ListOrdersWithoutCourier(ctx context.Context) ([]*Order, error)
	ListOrdersByCourierID(ctx context.Context, courierID int64) ([]*Order, error)
	GetOrder(ctx context.Context, id int64) (*Order, error)
	AcceptOrder(ctx context.Context, id int64, courierID int64) (*Order, error)
	CompleteOrder(ctx context.Context, id int64) (*Order, error)
}