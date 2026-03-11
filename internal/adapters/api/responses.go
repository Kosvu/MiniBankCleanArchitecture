package api

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type AddUserRequest struct {
	FullName string `json:"full_name"`
}

type GetUserRequest struct {
	ID uuid.UUID `json:"id"`
}

type TransactionRequest struct {
	Type   string `json:"type"`
	Amount int    `json:"amount"`
}

type ErrDTO struct {
	Error string    `json:"error"`
	Time  time.Time `json:"time"`
}

func (e *ErrDTO) ToString() string {
	b, err := json.MarshalIndent(e, "", "    ")
	if err != nil {
		panic(err)
	}

	return string(b)
}
