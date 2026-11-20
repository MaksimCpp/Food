package domain

import "errors"

var (
	ErrOrderNotFound          = errors.New("Order not found")
	ErrOrderAlreadyAccept     = errors.New("Order already accept")
	ErrInvalidOrderTransition = errors.New("Invalid orders transition")
)