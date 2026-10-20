package service

import (
	"context"
	"fmt"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/MaksimCpp/products/internal/dto"
	"github.com/google/uuid"
)

type PostgreSQLProductsService struct {
	repo    domain.ProductsRepository
	storage domain.FileStorage
}

func NewPostgreSQLProductsService(
	repo domain.ProductsRepository, storage domain.FileStorage,
) *PostgreSQLProductsService {
	return &PostgreSQLProductsService{
		repo: repo,
		storage: storage,
	}
}

func (s *PostgreSQLProductsService) Create(
	ctx context.Context, in *dto.CreateProductInput,
) (*dto.CreateProductOutput, error) {
	imageUuid := uuid.New().String()
	imageKey := fmt.Sprintf(
		"%s-%s",
		imageUuid,
		in.ImageName,
	)

	key, err := s.storage.Upload(ctx, imageKey, in.FileBytes, in.FileSize, in.FileType)

	product := domain.Product{
		Name: in.Name,
		Description: in.Description,
		Price: in.Price,
		ImageKey: key,
		IsAvailable: in.IsAvailable,
	}

	result, err := s.repo.Create(ctx, &product)

	if err != nil {
		return nil, err
	}

	out := dto.CreateProductOutput{
		ID: result.ID,
		Name: result.Name,
		Description: result.Description,
		Price: result.Price,
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &out, nil
}

func (s *PostgreSQLProductsService) Get(
	ctx context.Context, id int64,
) (*dto.GetProductOutput, error) {
	result, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	out := dto.GetProductOutput{
		Name: result.Name,
		Description: result.Description,
		Price: result.Price,
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &out, nil
}

func (s *PostgreSQLProductsService) List(
	ctx context.Context, limit int, offset int,
) ([]*dto.ListProductsOutput, error) {
	result, err := s.repo.List(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	out := make([]*dto.ListProductsOutput, 0, len(result))

	for _, product := range result {
		out = append(out, &dto.ListProductsOutput{
			Name: product.Name,
			Price: product.Price,
			ImageKey: product.ImageKey,
			IsAvailable: product.IsAvailable,
		})
	}

	return out, nil
}

func (s *PostgreSQLProductsService) Update(
	ctx context.Context, newProduct *dto.UpdateProductIntput, id int64,
) (*dto.UpdateProductOutput, error) {

	product := domain.Product{
		Name: newProduct.Name,
		Description: newProduct.Description,
		Price: newProduct.Price,
		ImageKey: newProduct.ImageKey,
		IsAvailable: newProduct.IsAvailable,
	}

	result, err := s.repo.Update(
		ctx, &product,id,
	)

	if err != nil {
		return nil, err
	}

	out := dto.UpdateProductOutput{
		ID: result.ID,
		Name: result.Name,
		Description: result.Description,
		Price: result.Price,
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &out, nil
}

func (s *PostgreSQLProductsService) Delete(ctx context.Context, id int64) error {
	return s.repo.Delete(ctx, id)
}