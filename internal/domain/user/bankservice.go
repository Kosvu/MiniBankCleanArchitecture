package domains

import (
	"context"

	"github.com/google/uuid"
)

type BankService interface {
	GetUserByID(ctx context.Context, id uuid.UUID) (User, error)
	GetAllUser(ctx context.Context, limit, offset int) ([]User, error)
	AddUser(ctx context.Context, fullName string) (User, error)
	CreateTransaction(ctx context.Context, class string, userID uuid.UUID, amount int) (User, error)
}

const (
	TxDeposit  = "deposit"
	TxWithdraw = "withdraw"
)

type bankService struct {
	storage Storage
}

func NewBankService(storage Storage) *bankService {
	return &bankService{storage: storage}
}

func (b *bankService) GetUserByID(ctx context.Context, id uuid.UUID) (User, error) {
	return b.storage.GetUser(id)
}

func (b *bankService) GetAllUser(ctx context.Context, limit, offset int) ([]User, error) {
	return b.storage.GetAllUsers(0, 0)
}

func (b *bankService) AddUser(ctx context.Context, fullName string) (User, error) {
	u, err := NewUser(fullName)
	if err != nil {
		return User{}, err
	}

	if err := b.storage.Create(u); err != nil {
		return User{}, err
	}

	return u, nil
}

func (b *bankService) CreateTransaction(ctx context.Context, class string, userID uuid.UUID, amount int) (User, error) {
	user, err := b.storage.GetUser(userID)
	if err != nil {
		return User{}, err
	}

	switch class {
	case TxDeposit:
		if err := user.Deposit(amount); err != nil {
			return User{}, err
		}
	case TxWithdraw:
		if err := user.Withdraw(amount); err != nil {
			return User{}, err
		}
	default:
		return User{}, err
	}
	if err := b.storage.Update(user); err != nil {
		return User{}, err
	}

	return user, nil
}
