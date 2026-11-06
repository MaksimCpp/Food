package domain

import (
	"context"

	"github.com/MaksimCpp/products/internal/dto"
)

type ProductsRepository interface {
	Create(ctx context.Context, product *Product) (*Product, error)
	Get(ctx context.Context, id int64) (*Product, error)
	List(ctx context.Context, limit int, offset int) ([]*Product, error)
	Update(ctx context.Context, newProduct *Product, id int64) (*Product, error)
	Delete(ctx context.Context, id int64) error
	GetIDs(ctx context.Context, limit int, offset int) ([]*Product, error)
	GetByIDs(ctx context.Context, ids ...int64) ([]*Product, error)
	Reserve(ctx context.Context, in *dto.ReserveProductInputOutput) (
		*dto.ReserveProductInputOutput, error,
	)
	DeductDeliveredStock(ctx context.Context, in *dto.DeductDeliveredStockInput) (
		*dto.DeductDeliveredStockOutput, error,
	)
}
