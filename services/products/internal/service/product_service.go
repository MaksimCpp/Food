package service

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/MaksimCpp/products/internal/dto"
	"github.com/google/uuid"
)

type PostgreSQLProductsService struct {
	repo    domain.ProductsRepository
	storage domain.FileStorage
	cache   domain.Cache
}

func NewPostgreSQLProductsService(
	repo domain.ProductsRepository, storage domain.FileStorage, cache domain.Cache,
) *PostgreSQLProductsService {
	return &PostgreSQLProductsService{
		repo: repo,
		storage: storage,
		cache: cache,
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

	dataForCache, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	cacheKey := fmt.Sprintf("product_%d", result.ID)

	err = s.cache.Set(ctx, cacheKey, dataForCache)

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
	key := fmt.Sprintf("product_%d", id)

	data, err := s.cache.Get(ctx, key)

	if err == nil {
		var product dto.GetProductOutput
		jsonErr := json.Unmarshal(data, &product)

		if jsonErr != nil {
			return nil, jsonErr
		}

		return &product, nil
	}

	result, err := s.repo.Get(ctx, id)

	if err != nil {
		return nil, err
	}

	data, err = json.Marshal(result)

	if err != nil {
		return nil, err
	}

	err = s.cache.Set(ctx, key, data)

	out := dto.GetProductOutput{
		Name: result.Name,
		Description: result.Description,
		Price: result.Price,
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &out, nil
}

// Сделать MGet из кеша
func (s *PostgreSQLProductsService) List(
	ctx context.Context, limit int, offset int,
) ([]*dto.ListProductsOutput, error) {
	ids, err := s.repo.GetIDs(ctx, limit, offset)

	if err != nil {
		return nil, err
	}

	keys := make([]string, len(ids))

	for i, k := range ids {
		keys[i] = fmt.Sprintf("product_%d", k.ID)
	}

	cacheResults, err := s.cache.MGet(ctx, keys...)

	if err != nil {
		result, pgErr := s.repo.List(ctx, limit, offset)

		if pgErr != nil {
			return nil, pgErr
		}

		out := make([]*dto.ListProductsOutput, 0, len(result))

		for _, product := range result {
			productOut := dto.ListProductsOutput{
				Name: product.Name,
				Price: product.Price,
				ImageKey: product.ImageKey,
				IsAvailable: product.IsAvailable,
			}

			out = append(out, &productOut)
		}
		return out, nil
	}

	missingsIDs := make([]int64, 0)

	for i, v := range cacheResults {
		if v == nil {
			missingsIDs = append(missingsIDs, ids[i].ID)
		}
	}

	if len(missingsIDs) > 0 {
		postgresResults, err := s.repo.GetByIDs(ctx, missingsIDs...)

		if err != nil {
			return nil, err
		}

		postresResultsMap := make(map[int64]*domain.Product)
	
		for _, product := range postgresResults {
			postresResultsMap[product.ID] = product
		}
	
		for i, res := range cacheResults {
			if res == nil {
				currentID := ids[i].ID
				product, ok := postresResultsMap[currentID]

				if !ok {
					continue
				}

				out := dto.ListProductsOutput{
					Name: product.Name,
					Price: product.Price,
					ImageKey: product.ImageKey,
					IsAvailable: product.IsAvailable,
				}

				JSONFromProduct, err := json.Marshal(out)
	
				if err != nil {
					return nil, err
				}
	
				cacheResults[i] = string(JSONFromProduct)

				_ = s.cache.Set(
					context.WithoutCancel(ctx), 
					fmt.Sprintf("product_%d", currentID), 
					JSONFromProduct,
				)
			}
		}
	}

	if err != nil {
		return nil, err
	}

	out := make([]*dto.ListProductsOutput, 0, len(cacheResults))

	for _, product := range cacheResults {
		var productFromJSON dto.ListProductsOutput

		strProduct, ok := product.(string)

		if !ok || strProduct == "" {
			continue
		}

		err = json.Unmarshal([]byte(strProduct), &productFromJSON)

		if err != nil {
			return nil, err
		}

		out = append(out, &productFromJSON)
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

	data, err := json.Marshal(result)

	if err != nil {
		return nil, err
	}

	key := fmt.Sprintf("product_%d", result.ID)

	err = s.cache.Set(ctx, key, data)

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
	err := s.repo.Delete(ctx, id)

	if err != nil {
		return nil
	}

	key := fmt.Sprintf("product_%d", id)

	return s.cache.Delete(ctx, key)
}