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
	DeleteUser(ctx context.Context, id uuid.UUID) error
	GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, error)
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
	return b.storage.GetUser(ctx, id)
}

func (b *bankService) GetUserTransactions(ctx context.Context, userID uuid.UUID) ([]Transaction, error) {
	return b.storage.GetUserTransactions(ctx, userID)
}

func (b *bankService) GetAllUser(ctx context.Context, limit, offset int) ([]User, error) {
	return b.storage.GetAllUsers(ctx, limit, offset)
}

func (b *bankService) AddUser(ctx context.Context, fullName string) (User, error) {
	u, err := NewUser(fullName)
	if err != nil {
		return User{}, err
	}

	if err := b.storage.Create(ctx, u); err != nil {
		return User{}, err
	}

	return u, nil
}

func (b *bankService) CreateTransaction(ctx context.Context, class string, userID uuid.UUID, amount int) (User, error) {
	user, err := b.storage.GetUser(ctx, userID)
	if err != nil {
		return User{}, ErrUserNotFound
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
		return User{}, ErrUnknownTransactionType
	}
	if err := b.storage.Update(ctx, user); err != nil {
		return User{}, err
	}

	tx := NewTransaction(userID, class, amount)

	err = b.storage.CreateTransactionRecord(ctx, tx)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

func (b *bankService) DeleteUser(ctx context.Context, id uuid.UUID) error {
	return b.storage.Delete(ctx, id)
}
