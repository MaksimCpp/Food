package service

import (
	"context"
	"fmt"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/MaksimCpp/products/internal/dto"
	"github.com/google/uuid"
)

// 	List(ctx context.Context, limit int, offset int) ([]*dto.ListProductsOutput, error)
// 	Update(ctx context.Context, newProduct *Product, id int64) (*dto.UpdateProductOutput, error)
// 	Delete(ctx context.Context, id int64) error

type PostgreSQLProductService struct {
	repo    domain.ProductRepository
	storage domain.FileStorage
}

func NewPostgreSQLProductService(
	repo domain.ProductRepository, storage domain.FileStorage,
) *PostgreSQLProductService {
	return &PostgreSQLProductService{
		repo: repo,
		storage: storage,
	}
}

func (s *PostgreSQLProductService) Create(
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

func (s *PostgreSQLProductService) Get(
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