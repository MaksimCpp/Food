package domain

import "context"

type ProductRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	Get(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, limit int, offset int) ([]*Product, error)
	Update(ctx context.Context, newProduct *Product, id int64) (*Product, error)
	Delete(ctx context.Context, id int64) (*Product, error)
}