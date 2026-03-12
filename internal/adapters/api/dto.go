package api

import (
	"encoding/json"
	"net/http"
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

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, err error) {
	writeJSON(w, status, ErrDTO{
		Error: err.Error(),
		Time:  time.Now(),
	})
}
