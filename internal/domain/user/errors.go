package domains

import "errors"

var (
	ErrUserNotFound           = errors.New("user not found")
	ErrEmptyFullName          = errors.New("full name is required")
	ErrInvalidAmount          = errors.New("amount must be positive")
	ErrInsufficientFunds      = errors.New("insufficient funds")
	ErrUnknownTransactionType = errors.New("unknown transaction type")
	ErrUserAlreadyExists      = errors.New("user already exists")
	ErrInvalidPagination      = errors.New("invalid pagination params")
)
