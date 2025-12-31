package shared

import "errors"

var (
	ErrNotFound           = errors.New("resource not found")
	ErrInvalidInput       = errors.New("invalid input")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrOrderNotModifiable = errors.New("order cannot be modified in current status")
)
