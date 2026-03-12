package domains

import (
	"github.com/google/uuid"
)

type User struct {
	ID       uuid.UUID `json:"id"`
	FullName string    `json:"full_name"`
	Balance  int       `json:"balance"`
}

func NewUser(fullName string) (User, error) {

	if fullName == "" {
		return User{}, ErrEmptyFullName
	}

	return User{
		ID:       uuid.New(),
		FullName: fullName,
		Balance:  0,
	}, nil
}

func (u *User) Deposit(amount int) error {
	if amount <= 0 {
		return ErrInvalidAmount
	}

	u.Balance += amount
	return nil
}

func (u *User) Withdraw(amount int) error {

	if amount <= 0 {
		return ErrInvalidAmount
	}

	if u.Balance < amount {
		return ErrInsufficientFunds
	}

	u.Balance -= amount
	return nil
}
