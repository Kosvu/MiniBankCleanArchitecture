package domains

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrEmptyFullName = errors.New("full_name is required")
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
		return errors.New("Deposit error!")
	}

	u.Balance += amount
	return nil
}

func (u *User) Withdraw(amount int) error {

	if amount <= 0 {
		return errors.New("Withdraw amount must be positive!")
	}

	if u.Balance < amount {
		return errors.New("Withdraw error!")
	}

	u.Balance -= amount
	return nil
}
