package domains

import (
	"errors"
	"testing"
)

func TestNewUser(t *testing.T) {
	t.Run("Succees", func(t *testing.T) {
		user, err := NewUser("Ivan Ivanov")

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
			t.Fatalf("expected generate ID")
		}
	})

	t.Run("Empty full name", func(t *testing.T) {
		_, err := NewUser("")

		if !errors.Is(err, ErrEmptyFullName) {
			t.Fatalf("expected ErrEmptyFullName, got %v", err)
		}
	})
}

func TestUserDeposit(t *testing.T) {
	t.Run("Succees", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Deposit(50)

		if err != nil {
			t.Fatal(err)
		}

		if u.Balance != 150 {
			t.Fatalf("expected balance 150, got %d", u.Balance)
		}
	})

	t.Run("Negative deposit", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Deposit(-100)

		if !errors.Is(err, ErrInvalidAmount) {
			t.Fatalf("expected ErrInvalidAmount, got %v", err)
		}

		if u.Balance != 100 {
			t.Fatalf("expected user balance 100, got %d", u.Balance)
		}
	})

	t.Run("Zero amount", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Deposit(0)

		if !errors.Is(err, ErrInvalidAmount) {
			t.Fatalf("expected ErrInvalidAmount, got %v", err)
		}

		if u.Balance != 100 {
			t.Fatalf("expected user balance 100, got %d", u.Balance)
		}
	})
}

func TestUserWithdraw(t *testing.T) {
	t.Run("Succees", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Withdraw(50)

		if err != nil {
			t.Fatal(err)
		}

		if u.Balance != 50 {
			t.Fatalf("expected balance 50, got %d", u.Balance)
		}
	})

	t.Run("zero amount", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Withdraw(0)

		if !errors.Is(err, ErrInvalidAmount) {
			t.Fatalf("expected ErrInvalidAmount, got %v", err)
		}

		if u.Balance != 100 {
			t.Fatalf("expected user balance 100, got %d", u.Balance)
		}
	})

	t.Run("neggative amount", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Withdraw(-100)

		if !errors.Is(err, ErrInvalidAmount) {
			t.Fatalf("expected ErrInavalidAmount, got %v", err)
		}

		if u.Balance != 100 {
			t.Fatalf("expected user balance 100, got %d", u.Balance)
		}
	})

	t.Run("Insufficient amount", func(t *testing.T) {
		u := User{
			FullName: "Ivan Ivanov",
			Balance:  100,
		}

		err := u.Withdraw(150)

		if !errors.Is(err, ErrInsufficientFunds) {
			t.Fatalf("expected ErrIsufficient, got %v", err)
		}

		if u.Balance != 100 {
			t.Fatalf("expected user balance 100, got %d", u.Balance)
		}
	})
}
