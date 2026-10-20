package domain

import (
	"context"

	"github.com/MaksimCpp/products/internal/dto"
)

type ProductsService interface {
	Create(ctx context.Context, in *dto.CreateProductInput) (*dto.CreateProductOutput, error)
	Get(ctx context.Context, id int64) (*dto.GetProductOutput, error)
	List(ctx context.Context, limit int, offset int) ([]*dto.ListProductsOutput, error)
	Update(
		ctx context.Context, newProduct *dto.UpdateProductIntput, id int64,
	) (*dto.UpdateProductOutput, error)
	Delete(ctx context.Context, id int64) error
}