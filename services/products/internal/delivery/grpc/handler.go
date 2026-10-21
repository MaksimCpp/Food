package deliverygrpc

import (
	"context"
	"errors"

	"github.com/MaksimCpp/products/internal/domain"
	"github.com/MaksimCpp/products/internal/dto"
	productspb "github.com/MaksimCpp/products/proto"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ProductsHandler struct {
	productspb.UnimplementedProductsServiceServer

	service domain.ProductsService
}

func NewProductsHandler(service domain.ProductsService) *ProductsHandler {
	return &ProductsHandler{
		service: service,
	}
}

func (h *ProductsHandler) CreateProduct(
	ctx context.Context, req *productspb.CreateProductRequest,
) (*productspb.CreateProductResponse, error) {
	in := dto.CreateProductInput{
		Name: req.Name,
		Description: req.Description,
		Price: int(req.Price),
		ImageName: req.ImageName,
		IsAvailable: req.IsAvailable,
	}

	result, err := h.service.Create(ctx, &in)

	if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error.")
	}

	res := productspb.CreateProductResponse{
		Id: result.ID,
		Name: result.Name,
		Description: result.Description,
		Price: int64(result.Price),
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &res, nil
}

func (h *ProductsHandler) GetProduct(
	ctx context.Context, req *productspb.GetProductRequest,
) (*productspb.GetProductResponse, error) {
	result, err := h.service.Get(ctx, req.Id)

	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "Product not found.")
		}
		return nil, status.Error(codes.Internal, "Internal server error.")
	}

	res := productspb.GetProductResponse{
		Id: result.ID,
		Name: result.Name,
		Description: result.Description,
		Price: int64(result.Price),
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &res, nil
}

func (h *ProductsHandler) ListProducts(
	ctx context.Context, req *productspb.ListProductsRequest,
) (*productspb.ListProductsResponse, error) {
	result, err := h.service.List(ctx, int(req.Limit), int(req.Offset))

	if err != nil {
		return nil, err
	}

	var res *productspb.ListProductsResponse

	for _, product := range result {
		res.Products = append(res.Products, &productspb.ProductItem{
			Name: product.Name,
			Price: int64(product.Price),
			ImageKey: product.ImageKey,
			IsAvailable: product.IsAvailable,
		})
	}

	return res, nil
}

func (h *ProductsHandler) UpdateProduct(
	ctx context.Context, req *productspb.UpdateProductRequest,
) (*productspb.UpdateProductResponse, error) {
	in := dto.UpdateProductIntput{
		ID: req.Id,
		Name: req.Name,
		Description: req.Description,
		Price: int(req.Price),
		ImageKey: req.ImageKey,
		IsAvailable: req.IsAvailable,
	}

	result, err := h.service.Update(ctx, &in, in.ID)

	if err != nil {
		if errors.Is(err, domain.ErrProductNotFound) {
			return nil, status.Error(codes.NotFound, "Product not found.")
		}
		return nil, status.Error(codes.Internal, "Internal server error.")
	}

	out := productspb.UpdateProductResponse{
		Id: result.ID,
		Name: result.Name,
		Description: result.Description,
		Price: int64(result.Price),
		ImageKey: result.ImageKey,
		IsAvailable: result.IsAvailable,
	}

	return &out, nil
}

func (h *ProductsHandler) DeleteProduct(
	ctx context.Context, req *productspb.DeleteProductRequest,
) (*emptypb.Empty, error) {
	return nil, h.service.Delete(ctx, req.Id)
}