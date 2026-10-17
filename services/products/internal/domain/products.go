package domain

type Product struct {
	ID          int64
	Name        string
	Description string
	Price       int
	ImageKey    string
	IsAvailable bool
}