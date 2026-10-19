package domain

import "errors"

var (
	ErrNoRowsAffected  = errors.New("No rows affected.")
	ErrProductNotFound = errors.New("Product not found")
)