package api

import (
	"errors"
	domains "minibank/internal/domain/users"
	"net/http"
)

func mapError(err error) int {
	switch {
	case errors.Is(err, domains.ErrUserNotFound):
		return http.StatusNotFound
	case errors.Is(err, domains.ErrEmptyFullName):
		return http.StatusBadRequest
	case errors.Is(err, domains.ErrInvalidAmount):
		return http.StatusBadRequest
	case errors.Is(err, domains.ErrUnknownTransactionType):
		return http.StatusBadRequest
	case errors.Is(err, domains.ErrInsufficientFunds):
		return http.StatusConflict
	case errors.Is(err, domains.ErrInvalidPagination):
		return http.StatusBadRequest
	default:
		return http.StatusInternalServerError
	}
}
