package dto

type CreateProductInput struct {
	Name          string
	Description   string
	Price         int
	ImageName     string
	IsAvailable   bool
	StockQuantity int64

	FileBytes []byte
	FileSize  int64
	FileType  string
}

type CreateProductOutput struct {
	ID            int64
	Name          string
	Description   string
	Price         int
	ImageKey      string
	IsAvailable   bool
	StockQuantity int64
}

type GetProductOutput struct {
	ID            int64
	Name          string
	Description   string
	Price         int
	ImageKey      string
	IsAvailable   bool
	StockQuantity int64
}

type UpdateProductIntput struct {
	ID            int64
	Name          string
	Description   string
	Price         int
	ImageKey      string
	IsAvailable   bool
	StockQuantity int64
}

type UpdateProductOutput struct {
	ID            int64
	Name          string
	Description   string
	Price         int
	ImageKey      string
	IsAvailable   bool
	StockQuantity int64
}

type ReserveProductInputOutput struct {
	ID       int64
	Quantity int64
}

type DeductDeliveredStockInput struct {
	ID       int64
	Quantity int64
}

type DeductDeliveredStockOutput struct {
	ID       int64
	StockQuantity int64
}
