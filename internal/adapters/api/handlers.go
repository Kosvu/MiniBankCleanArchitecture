package api

import (
	"encoding/json"
	domains "minibank/internal/domain/user"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type HTTPHandlers struct {
	bank domains.BankService
}

func NewHTTPHandlers(bank domains.BankService) *HTTPHandlers {
	return &HTTPHandlers{
		bank: bank,
	}
}

func (h *HTTPHandlers) GetAllUsersH(w http.ResponseWriter, r *http.Request) {
	page := 1
	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		pageInt, err := strconv.Atoi(pageStr)
		if err != nil {
			writeError(w, mapError(domains.ErrInvalidPagination), err)
			return
		}

		if pageInt < 1 {
			writeError(w, mapError(domains.ErrInvalidPagination), err)
			return
		}
		page = pageInt
	}
	limit := 3
	offset := (page - 1) * limit

	users, err := h.bank.GetAllUser(r.Context(), limit, offset)
	if err != nil {
		writeError(w, mapError(err), err)
		return
	}
	writeJSON(w, http.StatusOK, users)
}

func (h *HTTPHandlers) GetUserByIDH(w http.ResponseWriter, r *http.Request) {

	idStr := mux.Vars(r)["id"]

	id, err := uuid.Parse(idStr) // перевод id в UUID
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.bank.GetUserByID(r.Context(), id)

	if err != nil {
		writeError(w, mapError(err), err)
		return
	}

	writeJSON(w, http.StatusOK, user)
}

func (h *HTTPHandlers) AddUserH(w http.ResponseWriter, r *http.Request) {
	var userDTO AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.bank.AddUser(r.Context(), userDTO.FullName)

	if err != nil {
		writeError(w, mapError(err), err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *HTTPHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]

	id, err := uuid.Parse(idString)
	if err != nil {
		writeError(w, mapError(err), err)
		return
	}

	var transaction TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	user, err := h.bank.CreateTransaction(r.Context(), transaction.Type, id, transaction.Amount)
	if err != nil {
		writeError(w, mapError(err), err)
		return
	}

	writeJSON(w, http.StatusCreated, user)
}

func (h *HTTPHandlers) DeleteUserH(w http.ResponseWriter, r *http.Request) {
	idString := mux.Vars(r)["id"]

	idUuid, err := uuid.Parse(idString)
	if err != nil {
		writeError(w, http.StatusBadRequest, err)
		return
	}

	if err := h.bank.DeleteUser(r.Context(), idUuid); err != nil {
		writeError(w, mapError(err), err)
		return
	}

	w.WriteHeader(http.StatusOK)
}
