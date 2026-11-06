package domain

import "errors"

var (
	ErrProductNotFound = errors.New("Product not found")
	ErrProductQuantity = errors.New("Invalid quantity")
)
