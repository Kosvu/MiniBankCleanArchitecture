package domains

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
)

func TestBankService_AddUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		user, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		if user.FullName != "Ivan Ivanov" {
			t.Fatalf("expected full name %q, got %q", "Ivan Ivanov", user.FullName)
		}

		if user.Balance != 0 {
			t.Fatalf("expected balance %d, got %d", 0, user.Balance)
		}

		if user.ID.String() == "" {
			t.Fatal("expected generated id")
		}
	})

	t.Run("empty full name", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		_, err := service.AddUser(ctx, "")
		if !errors.Is(err, ErrEmptyFullName) {
			t.Fatalf("expected ErrEmptyFullName, got %v", err)
		}
	})
}

func TestBankService_GetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		createdUser, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		gotUser, err := service.GetUserByID(ctx, createdUser.ID)
		if err != nil {
			t.Fatal(err)
		}

		if gotUser.ID != createdUser.ID {
			t.Fatalf("expected id %v, got %v", createdUser.ID, gotUser.ID)
		}

		if gotUser.FullName != createdUser.FullName {
			t.Fatalf("expected full name %q, got %q", createdUser.FullName, gotUser.FullName)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		_, err := service.GetUserByID(ctx, uuid.New())
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestBankService_CreateTransaction(t *testing.T) {
	t.Run("deposit success", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		user, err := service.AddUser(ctx, "Ivan")
		if err != nil {
			t.Fatal(err)
		}

		updatedUser, err := service.CreateTransaction(ctx, TxDeposit, user.ID, 100)
		if err != nil {
			t.Fatal(err)
		}

		if updatedUser.Balance != 100 {
			t.Fatalf("expected balance %d, got %d", 100, updatedUser.Balance)
		}
	})

	t.Run("withdraw success", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		user, err := service.AddUser(ctx, "Ivan")
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.CreateTransaction(ctx, TxDeposit, user.ID, 100)
		if err != nil {
			t.Fatal(err)
		}

		updatedUser, err := service.CreateTransaction(ctx, TxWithdraw, user.ID, 40)
		if err != nil {
			t.Fatal(err)
		}

		if updatedUser.Balance != 60 {
			t.Fatalf("expected balance %d, got %d", 60, updatedUser.Balance)
		}
	})

	t.Run("insufficient funds", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		user, err := service.AddUser(ctx, "Ivan")
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.CreateTransaction(ctx, TxWithdraw, user.ID, 100)
		if !errors.Is(err, ErrInsufficientFunds) {
			t.Fatalf("expected ErrInsufficientFunds, got %v", err)
		}
	})

	t.Run("unknown transaction type", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		user, err := service.AddUser(ctx, "Ivan")
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.CreateTransaction(ctx, "unknown", user.ID, 100)
		if !errors.Is(err, ErrUnknownTransactionType) {
			t.Fatalf("expected ErrUnknownTransactionType, got %v", err)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		_, err := service.CreateTransaction(ctx, TxDeposit, uuid.New(), 100)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestBankService_DeleteUser(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storage := NewStorage()
		service := NewBankService(storage)
		ctx := context.Background()

		user, err := service.AddUser(ctx, "Ivan")
		if err != nil {
			t.Fatal(err)
		}

		err = service.DeleteUser(ctx, user.ID)
		if err != nil {
			t.Fatal(err)
		}

		_, err = service.GetUserByID(ctx, user.ID)
		if !errors.Is(err, ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound after delete, got %v", err)
		}
	})
}
