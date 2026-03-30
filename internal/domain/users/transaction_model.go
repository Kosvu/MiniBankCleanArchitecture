package domains

import (
	"time"

	"github.com/google/uuid"
)

type Transaction struct {
	ID              uuid.UUID `json:"id"`
	UserID          uuid.UUID `json:"user_id"`
	TransactionType string    `json:"type"`
	Amount          int       `json:"amount"`
	CreatedAt       time.Time `json:"created_at"`
}

func NewTransaction(userID uuid.UUID, txType string, amount int) Transaction {
	return Transaction{
		ID:              uuid.New(),
		UserID:          userID,
		TransactionType: txType,
		Amount:          amount,
		CreatedAt:       time.Now(),
	}
}
