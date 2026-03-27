package api

import (
	"context"
	"log/slog"
	domains "minibank/internal/domain/users"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func newTestHandlers() *HTTPHandlers {
	storage := domains.NewStorage()
	service := domains.NewBankService(storage)
	logger := slog.Default()

	return NewHTTPHandlers(service, logger)
}

func TestHealthCheckH(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := newTestHandlers()

		req := httptest.NewRequest(http.MethodGet, "/bank/health", nil)
		rec := httptest.NewRecorder()

		h.HealthCheckH(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
		}

		if strings.TrimSpace(rec.Body.String()) != "OK" {
			t.Fatalf("expected %q, got %q", "OK", rec.Body.String())
		}
	})
}

func TestAddUserH(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		h := newTestHandlers()

		body := `{"full_name": "Ivan Ivanov"}`
		req := httptest.NewRequest(http.MethodPost, "/bank", strings.NewReader(body))
		rec := httptest.NewRecorder()

		h.AddUserH(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected %d, got %d", http.StatusCreated, rec.Code)
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		h := newTestHandlers()

		body := `{"full_name":}`
		req := httptest.NewRequest(http.MethodPost, "/bank", strings.NewReader(body))
		rec := httptest.NewRecorder()

		h.AddUserH(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})

	t.Run("empty full name", func(t *testing.T) {
		h := newTestHandlers()

		body := `{"full_name": ""}`
		req := httptest.NewRequest(http.MethodPost, "/bank", strings.NewReader(body))
		rec := httptest.NewRecorder()

		h.AddUserH(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})
}

func TestGetUserByID(t *testing.T) {
	storage := domains.NewStorage()
	service := domains.NewBankService(storage)
	logger := slog.Default()

	h := NewHTTPHandlers(service, logger)
	rout := NewRouter(h)

	ctx := context.Background()
	u, err := service.AddUser(ctx, "Ivan Ivanov")
	if err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodGet, "/bank/"+u.ID.String(), nil)
	rec := httptest.NewRecorder()

	rout.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, rec.Code)
	}

	t.Run("invalid json", func(t *testing.T) {
		h := newTestHandlers()
		rout := NewRouter(h)

		req := httptest.NewRequest(http.MethodGet, "/bank/not-uuid-string", nil)
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d", http.StatusBadRequest, rec.Code)
		}
	})
}

func TestGetAllUsersH(t *testing.T) {
	t.Run("without page", func(t *testing.T) {
		h := newTestHandlers()

		req := httptest.NewRequest(http.MethodGet, "/bank", nil)
		rec := httptest.NewRecorder()

		h.GetAllUsersH(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("valid page", func(t *testing.T) {
		h := newTestHandlers()

		req := httptest.NewRequest(http.MethodGet, "/bank?page=2", nil)
		rec := httptest.NewRecorder()

		h.GetAllUsersH(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("non numeric page", func(t *testing.T) {
		h := newTestHandlers()

		req := httptest.NewRequest(http.MethodGet, "/bank?page=abc", nil)
		rec := httptest.NewRecorder()

		h.GetAllUsersH(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})

	t.Run("negative page", func(t *testing.T) {
		h := newTestHandlers()

		req := httptest.NewRequest(http.MethodGet, "/bank?page=-2", nil)
		rec := httptest.NewRecorder()

		h.GetAllUsersH(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})
}

func TestCreateTransactionH(t *testing.T) {
	t.Run("success deposit", func(t *testing.T) {
		storage := domains.NewStorage()
		service := domains.NewBankService(storage)
		logger := slog.Default()

		h := NewHTTPHandlers(service, logger)
		rout := NewRouter(h)

		ctx := context.Background()
		u, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		body := `{
		"type": "deposit",
		"amount": 100
		}`

		req := httptest.NewRequest(http.MethodPost, "/bank/"+u.ID.String()+"/transaction", strings.NewReader(body))
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusCreated, rec.Code, rec.Body.String())
		}

	})

	t.Run("success withdraw", func(t *testing.T) {
		storage := domains.NewStorage()
		service := domains.NewBankService(storage)
		logger := slog.Default()

		h := NewHTTPHandlers(service, logger)
		rout := NewRouter(h)

		ctx := context.Background()
		u, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}
		_, err = service.CreateTransaction(ctx, "deposit", u.ID, 150)
		if err != nil {
			t.Fatal(err)
		}

		body := `{
		"type": "withdraw",
		"amount": 100
		}`

		req := httptest.NewRequest(http.MethodPost, "/bank/"+u.ID.String()+"/transaction", strings.NewReader(body))
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusCreated {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusCreated, rec.Code, rec.Body.String())
		}
	})

	t.Run("invalid json", func(t *testing.T) {
		storage := domains.NewStorage()
		service := domains.NewBankService(storage)
		logger := slog.Default()

		h := NewHTTPHandlers(service, logger)
		rout := NewRouter(h)

		ctx := context.Background()
		u, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		body := `{
		"type": "withdraw",
		"amount": "100"
		}`

		req := httptest.NewRequest(http.MethodPost, "/bank/"+u.ID.String()+"/transaction", strings.NewReader(body))
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		h := newTestHandlers()
		rout := NewRouter(h)

		body := `{
		"type": "withdraw",
		"amount": "100"
		}`

		req := httptest.NewRequest(http.MethodPost, "/bank/not-uuid-string/transaction", strings.NewReader(body))
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})

	t.Run("unknown type", func(t *testing.T) {
		storage := domains.NewStorage()
		service := domains.NewBankService(storage)
		logger := slog.Default()

		h := NewHTTPHandlers(service, logger)
		rout := NewRouter(h)

		ctx := context.Background()
		u, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		body := `{
		"type": "non-transaction-type",
		"amount": 100
		}`

		req := httptest.NewRequest(http.MethodPost, "/bank/"+u.ID.String()+"/transaction", strings.NewReader(body))
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})

	t.Run("insufficient funds", func(t *testing.T) {
		storage := domains.NewStorage()
		service := domains.NewBankService(storage)
		logger := slog.Default()

		h := NewHTTPHandlers(service, logger)
		rout := NewRouter(h)

		ctx := context.Background()
		u, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}
		_, err = service.CreateTransaction(ctx, "deposit", u.ID, 100)
		if err != nil {
			t.Fatal(err)
		}

		body := `{
		"type": "withdraw",
		"amount": 150
		}`

		req := httptest.NewRequest(http.MethodPost, "/bank/"+u.ID.String()+"/transaction", strings.NewReader(body))
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusConflict {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusConflict, rec.Code, rec.Body.String())
		}
	})
}

func TestDeleteH(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		storage := domains.NewStorage()
		service := domains.NewBankService(storage)
		logger := slog.Default()
		ctx := context.Background()
		h := NewHTTPHandlers(service, logger)

		rout := NewRouter(h)

		u, err := service.AddUser(ctx, "Ivan Ivanov")
		if err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodDelete, "/bank/"+u.ID.String(), nil)
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusOK {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusOK, rec.Code, rec.Body.String())
		}
	})

	t.Run("invalid uuid", func(t *testing.T) {
		h := newTestHandlers()
		rout := NewRouter(h)

		req := httptest.NewRequest(http.MethodDelete, "/bank/not-uuid-string", nil)
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusBadRequest {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusBadRequest, rec.Code, rec.Body.String())
		}
	})

	t.Run("not found user", func(t *testing.T) {
		h := newTestHandlers()
		rout := NewRouter(h)
		randomID := uuid.New()

		req := httptest.NewRequest(http.MethodDelete, "/bank/"+randomID.String(), nil)
		rec := httptest.NewRecorder()

		rout.ServeHTTP(rec, req)

		if rec.Code != http.StatusNotFound {
			t.Fatalf("expected %d, got %d, body=%q", http.StatusNotFound, rec.Code, rec.Body.String())
		}
	})
}
