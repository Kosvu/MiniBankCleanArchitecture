package api

import (
	"context"
	"encoding/json"
	"log/slog"
	"minibank/internal/adapters/db"
	domains "minibank/internal/domain/users"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/joho/godotenv"
)

func newE2ERepo(t *testing.T) *db.UserRepository {
	t.Helper()

	if err := godotenv.Load("../../../.env.test"); err != nil {
		t.Fatalf("failed to load .env.test: %v", err)
	}

	ctx := context.Background()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		t.Fatal("TEST_DATABASE_URL is not set")
	}

	pool, err := db.NewConnection(ctx, dsn)
	if err != nil {
		t.Fatal(err)
	}

	t.Cleanup(func() {
		pool.Close()
	})

	repo := db.NewUserRepository(pool, slog.Default())

	_, err = pool.Exec(ctx, `TRUNCATE TABLE transactions, users CASCADE`)
	if err != nil {
		t.Fatal(err)
	}

	return repo
}

func TestE2E_UserLifecycle(t *testing.T) {
	repo := newE2ERepo(t)
	service := domains.NewBankService(repo)
	logger := slog.Default()
	h := NewHTTPHandlers(service, logger)
	router := NewRouter(h)

	var createdUser domains.User

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}
	ctx := context.Background()

	t.Run("create and capture user", func(t *testing.T) {
		body := `{"full_name":"Pavel Pavlov"}`

		resp, err := client.Post(server.URL+"/bank", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
			t.Fatal(err)
		}

		if createdUser.FullName != "Pavel Pavlov" {
			t.Fatalf("expected %q, got %q", "Pavel Pavlov", createdUser.FullName)
		}

		if createdUser.Balance != 0 {
			t.Fatalf("expected %d, got %d", 0, createdUser.Balance)
		}
	})

	t.Run("deposit user", func(t *testing.T) {
		body := `{
			"type": "deposit",
			"amount": 100
		}`

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			server.URL+"/bank/"+createdUser.ID.String()+"/transaction",
			strings.NewReader(body),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var updatedUser domains.User
		if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
			t.Fatal(err)
		}

		if updatedUser.Balance != 100 {
			t.Fatalf("expected balance %d, got %d", 100, updatedUser.Balance)
		}
	})

	t.Run("withdraw user", func(t *testing.T) {
		body := `{
			"type": "withdraw",
			"amount": 50
		}`

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			server.URL+"/bank/"+createdUser.ID.String()+"/transaction",
			strings.NewReader(body),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		var updatedUser domains.User
		if err := json.NewDecoder(resp.Body).Decode(&updatedUser); err != nil {
			t.Fatal(err)
		}

		if updatedUser.Balance != 50 {
			t.Fatalf("expected balance %d, got %d", 50, updatedUser.Balance)
		}
	})

	t.Run("get user after transaction", func(t *testing.T) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			server.URL+"/bank/"+createdUser.ID.String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var gotUser domains.User
		if err := json.NewDecoder(resp.Body).Decode(&gotUser); err != nil {
			t.Fatal(err)
		}

		if gotUser.ID != createdUser.ID {
			t.Fatalf("expected id %v, got %v", createdUser.ID, gotUser.ID)
		}

		if gotUser.FullName != "Pavel Pavlov" {
			t.Fatalf("expected FullName=%q, got %q", "Pavel Pavlov", gotUser.FullName)
		}

		if gotUser.Balance != 50 {
			t.Fatalf("expected balance %d, got %d", 50, gotUser.Balance)
		}
	})

	t.Run("delete user", func(t *testing.T) {

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodDelete,
			server.URL+"/bank/"+createdUser.ID.String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}
	})

	t.Run("get deleted user", func(t *testing.T) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			server.URL+"/bank/"+createdUser.ID.String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusNotFound {
			t.Fatalf("expected %d, got %d", http.StatusNotFound, resp.StatusCode)
		}
	})
}

func TestE2E_InsufficientFunds(t *testing.T) {

	repo := newE2ERepo(t)
	service := domains.NewBankService(repo)
	logger := slog.Default()
	h := NewHTTPHandlers(service, logger)
	router := NewRouter(h)

	var createdUser domains.User

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}
	ctx := context.Background()

	t.Run("create and capture user", func(t *testing.T) {
		body := `{"full_name":"Pavel Pavlov"}`

		resp, err := client.Post(server.URL+"/bank", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
			t.Fatal(err)
		}

		if createdUser.FullName != "Pavel Pavlov" {
			t.Fatalf("expected %q, got %q", "Pavel Pavlov", createdUser.FullName)
		}

		if createdUser.Balance != 0 {
			t.Fatalf("expected %d, got %d", 0, createdUser.Balance)
		}
	})

	t.Run("withdraw with InsufficientFunds", func(t *testing.T) {
		body := `{
			"type": "withdraw",
			"amount": 50
		}`

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			server.URL+"/bank/"+createdUser.ID.String()+"/transaction",
			strings.NewReader(body),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusConflict {
			t.Fatalf("expected %d, got %d", http.StatusConflict, resp.StatusCode)
		}
	})

	t.Run("balance remains zero", func(t *testing.T) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			server.URL+"/bank/"+createdUser.ID.String(),
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var gotUser domains.User

		if err := json.NewDecoder(resp.Body).Decode(&gotUser); err != nil {
			t.Fatal(err)
		}

		if gotUser.ID != createdUser.ID {
			t.Fatalf("expected id %v, got %v", createdUser.ID, gotUser.ID)
		}

		if gotUser.FullName != createdUser.FullName {
			t.Fatalf("expected %q, got %q", createdUser.FullName, gotUser.FullName)
		}

		if gotUser.Balance != 0 {
			t.Fatalf("expected %d balance, got %d", 0, gotUser.Balance)
		}
	})
}

func TestE2E_GetUserTransactions(t *testing.T) {
	repo := newE2ERepo(t)
	service := domains.NewBankService(repo)
	logger := slog.Default()
	h := NewHTTPHandlers(service, logger)
	router := NewRouter(h)

	server := httptest.NewServer(router)
	defer server.Close()

	client := &http.Client{}
	ctx := context.Background()

	var createdUser domains.User

	t.Run("create and capture user", func(t *testing.T) {
		body := `{"full_name":"Ivan Ivanov"}`

		resp, err := client.Post(server.URL+"/bank", "application/json", strings.NewReader(body))
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}

		if err := json.NewDecoder(resp.Body).Decode(&createdUser); err != nil {
			t.Fatal(err)
		}

		if createdUser.FullName != "Ivan Ivanov" {
			t.Fatalf("expected %q, got %q", "Ivan Ivanov", createdUser.FullName)
		}

		if createdUser.Balance != 0 {
			t.Fatalf("expected balance %d, got %d", 0, createdUser.Balance)
		}
	})

	t.Run("deposit", func(t *testing.T) {
		body := `{
			"type": "deposit",
			"amount": 100
		}`

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			server.URL+"/bank/"+createdUser.ID.String()+"/transaction",
			strings.NewReader(body),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})

	t.Run("withdraw", func(t *testing.T) {
		body := `{
			"type": "withdraw",
			"amount": 40
		}`

		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodPost,
			server.URL+"/bank/"+createdUser.ID.String()+"/transaction",
			strings.NewReader(body),
		)
		if err != nil {
			t.Fatal(err)
		}

		req.Header.Set("Content-Type", "application/json")

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, resp.StatusCode)
		}
	})

	t.Run("get user transactions", func(t *testing.T) {
		req, err := http.NewRequestWithContext(
			ctx,
			http.MethodGet,
			server.URL+"/bank/"+createdUser.ID.String()+"/transaction",
			nil,
		)
		if err != nil {
			t.Fatal(err)
		}

		resp, err := client.Do(req)
		if err != nil {
			t.Fatal(err)
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, resp.StatusCode)
		}

		var transactions []domains.Transaction
		if err := json.NewDecoder(resp.Body).Decode(&transactions); err != nil {
			t.Fatal(err)
		}

		if len(transactions) != 2 {
			t.Fatalf("expected %d transactions, got %d", 2, len(transactions))
		}
	})
}
