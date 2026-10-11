package domain

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       int
	ImageUrl    string
	IsAvailable bool
}