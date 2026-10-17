package dto

type CreateProductInput struct {
	Name        string
	Description string
	Price       int
	ImageName   string
	IsAvailable bool
}

type CreateProductOutput struct {
	ID          int64
	Name        string
	Description string
	Price       int
	ImageKey    string
	IsAvailable bool
}

type GetProductOutput struct {
	ID          int64
	Name        string
	Description string
	Price       int
	ImageKey    string
	IsAvailable bool
}

type ListProductsOutput struct {
	Name        string
	Price       int
	ImageKey    string
	IsAvailable bool
}

type UpdateProductOutput struct {
	ID          int64
	Name        string
	Description string
	Price       int
	ImageKey    string
	IsAvailable bool
}
