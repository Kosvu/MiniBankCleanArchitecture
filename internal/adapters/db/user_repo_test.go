package db

import (
	"context"
	"errors"
	"log/slog"
	domains "minibank/internal/domain/users"
	"os"
	"testing"

	"github.com/google/uuid"
	"github.com/joho/godotenv"
)

func NewTestRepo(t *testing.T) *UserRepository {
	t.Helper()

	if err := godotenv.Load("../../../.env.test"); err != nil {
		t.Fatalf("failed to load .env.test: %v", err)
	}

	ctx := context.Background()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	pool, err := NewConnection(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	repo := NewUserRepository(pool, slog.Default())

	_, err = pool.Exec(ctx, `TRUNCATE TABLE users`)

	if err != nil {
		t.Fatal(err)
	}

	return repo
}

func TestCreate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		user, err := domains.NewUser("Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		err = r.Create(ctx, user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		gotUser, err := r.GetUser(ctx, user.ID)
		if err != nil {
			t.Fatalf("failed to get created user: %v", err)
		}

		if gotUser.ID != user.ID {
			t.Fatalf("expected id %v, got %v", user.ID, gotUser.ID)
		}

		if gotUser.FullName != user.FullName {
			t.Fatalf("expected FullName %q, got %q", user.FullName, gotUser.FullName)
		}

		if gotUser.Balance != user.Balance {
			t.Fatalf("expected balance %d, got %d", user.Balance, gotUser.Balance)
		}
	})
}

func TestGetUserByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		user, err := domains.NewUser("Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		err = r.Create(ctx, user)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		gotUser, err := r.GetUser(ctx, user.ID)
		if err != nil {
			t.Fatalf("failed to get created user: %v", err)
		}

		if gotUser.ID != user.ID {
			t.Fatalf("expected id %v, got %v", user.ID, gotUser.ID)
		}

		if gotUser.FullName != user.FullName {
			t.Fatalf("expected FullName %q, got %q", user.FullName, gotUser.FullName)
		}

		if gotUser.Balance != user.Balance {
			t.Fatalf("expected balance %d, got %d", user.Balance, gotUser.Balance)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		_, err := r.GetUser(ctx, uuid.New())
		if !errors.Is(err, domains.ErrUserNotFound) {
			t.Fatalf("expected %v, got %v", domains.ErrUserNotFound, err)
		}
	})
}

func TestGetAllUsers(t *testing.T) {
	t.Run("empty result", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		users, err := r.GetAllUsers(ctx, 10, 0)
		if err != nil {
			t.Fatal(err)
		}

		if len(users) != 0 {
			t.Fatalf("expected 0 users, got %d", len(users))
		}
	})

	t.Run("several users", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		u1, err := domains.NewUser("Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}
		u2, err := domains.NewUser("Pavel Pavlov")
		if err != nil {
			t.Fatal(err)
		}

		if err := r.Create(ctx, u1); err != nil {
			t.Fatal(err)
		}

		if err := r.Create(ctx, u2); err != nil {
			t.Fatal(err)
		}

		users, err := r.GetAllUsers(ctx, 10, 0)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		if len(users) != 2 {
			t.Fatalf("expected 2 users, got %d", len(users))
		}
	})
}

func TestUpdate(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		u, err := domains.NewUser("Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Create(ctx, u)
		if err != nil {
			t.Fatal(err)
		}

		u.FullName = "Egor Egorov"
		u.Balance = 500

		err = r.Update(ctx, u)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		updatedUser, err := r.GetUser(ctx, u.ID)
		if err != nil {
			t.Fatal(err)
		}

		if updatedUser.Balance != 500 {
			t.Fatalf("expected balance 500, got %d", updatedUser.Balance)
		}

		if updatedUser.FullName != "Egor Egorov" {
			t.Fatalf("expected full name %q, got %q", "Egor Egorov", updatedUser.FullName)
		}
	})

	t.Run("user not found", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		u, err := domains.NewUser("Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		err = r.Update(ctx, u)
		if !errors.Is(err, domains.ErrUserNotFound) {
			t.Fatalf("expected ErrUserNotFound, got %v", err)
		}
	})
}

func TestDelete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		u, err := domains.NewUser("Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}
		err = r.Create(ctx, u)
		if err != nil {
			t.Fatal(err)
		}

		err = r.Delete(ctx, u.ID)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		_, err = r.GetUser(ctx, u.ID)
		if !errors.Is(err, domains.ErrUserNotFound) {
			t.Fatal("expected error after delete, got nil")
		}
	})

	t.Run("user not found", func(t *testing.T) {
		r := NewTestRepo(t)
		ctx := context.Background()

		err := r.Delete(ctx, uuid.New())
		if !errors.Is(err, domains.ErrUserNotFound) {
			t.Fatalf("expected user not found, got %v", err)
		}
	})
}
