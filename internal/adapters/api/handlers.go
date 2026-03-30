package api

import (
	"encoding/json"
	"log/slog"
	domains "minibank/internal/domain/users"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	bank domains.BankService
	log  *slog.Logger
}

func NewHTTPHandlers(bank domains.BankService, log *slog.Logger) *HTTPHandlers {
	return &HTTPHandlers{
		bank: bank,
		log:  log,
	}
}

func (h *HTTPHandlers) GetAllUsersH(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			h.log.Warn("invalid page query parameter", "page", pageStr, "err", err)
			writeError(w, mapError(domains.ErrInvalidPagination), err)
			return
		}

		if pageInt < 1 {
			h.log.Warn("page must be greater than zero", "page", pageInt)
			writeError(w, mapError(domains.ErrInvalidPagination), domains.ErrInvalidPagination)
			return
		}
		page = pageInt
	}
	limit := 3
	offset := (page - 1) * limit

	users, err := h.bank.GetAllUser(r.Context(), limit, offset)
	if err != nil {
		h.log.Error("failed to get all users", "err", err, "limit", limit, "offset", offset)
		writeError(w, mapError(err), err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *HTTPHandlers) GetUserByIDH(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr) // перевод id в UUID
	if err != nil {
		h.log.Warn("invalid user id", "id", idStr, "err", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.bank.GetUserByID(r.Context(), id)

	if err != nil {
		h.log.Error("failed to get user by id", "err", err, "user_id", id)
		writeError(w, mapError(err), err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *HTTPHandlers) AddUserH(w http.ResponseWriter, r *http.Request) {
	var userDTO AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		h.log.Warn("failed to decode add user request body", "err", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.bank.AddUser(r.Context(), userDTO.FullName)

	if err != nil {
		h.log.Error("failed to add user", "err", err, "full_name", userDTO.FullName)
		writeError(w, mapError(err), err)
		return
	}

	h.log.Info("user created", "user_id", user.ID)
	writeJSON(w, http.StatusCreated, user)
}

func (h *HTTPHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]

	id, err := uuid.Parse(idString)
	if err != nil {
		h.log.Warn("invalid user id for transaction", "id", idString, "err", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	var transaction TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		h.log.Warn("failed to decode transaction request body", "err", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.bank.CreateTransaction(r.Context(), transaction.Type, id, transaction.Amount)
	if err != nil {
		h.log.Error(
			"failed to create transaction",
			"err", err,
			"user_id", id,
			"type", transaction.Type,
			"amount", transaction.Amount,
		)
		writeError(w, mapError(err), err)
		return
	}

	h.log.Info(
		"transaction created",
		"user_id", id,
		"type", transaction.Type,
		"amount", transaction.Amount,
	)
	writeJSON(w, http.StatusCreated, user)
}

func (h *HTTPHandlers) DeleteUserH(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]

	idUuid, err := uuid.Parse(idString)
	if err != nil {
		h.log.Warn("invalid user id for delete", "id", idString, "err", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.bank.DeleteUser(r.Context(), idUuid); err != nil {
		h.log.Error("failed to delete user", "err", err, "user_id", idUuid)
		writeError(w, mapError(err), err)
		return
	}

	h.log.Info("user deleted", "user_id", idUuid)
	w.WriteHeader(http.StatusOK)
}

func (h *HTTPHandlers) GetUserTransactionsH(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]

	id, err := uuid.Parse(idString)
	if err != nil {
		h.log.Warn("invalid user id for get transactions", "id", idString, "err", err)
		writeError(w, http.StatusBadRequest, err)
		return
	}

	transactions, err := h.bank.GetUserTransactions(r.Context(), id)
	if err != nil {
		h.log.Error("failed to get user transactions", "err", err, "user_id", id)
		writeError(w, mapError(err), err)
		return
	}

	writeJSON(w, http.StatusOK, transactions)
}

func (h *HTTPHandlers) HealthCheckH(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("OK"))
}
