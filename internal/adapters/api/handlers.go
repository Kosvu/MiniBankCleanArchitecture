package api

import (
	"encoding/json"
	domains "minibank/internal/domain/user"
	"net/http"
	"time"

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
	users, err := h.bank.GetAllUser(r.Context(), 0, 0)
	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(users, "", "   ")
	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(b); err != nil {
		return
	}
}

func (h *HTTPHandlers) GetUserByIDH(w http.ResponseWriter, r *http.Request) {

	idSTR := mux.Vars(r)["id"]

	id, err := uuid.Parse(idSTR) // перевод id в UUID
	if err != nil {
		http.Error(w, "invaild uuid", http.StatusBadRequest)
	}

	user, err := h.bank.GetUserByID(r.Context(), id)

	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	b, err := json.MarshalIndent(user, "", "    ")

	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	if _, err := w.Write(b); err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}
}

func (h *HTTPHandlers) AddUserH(w http.ResponseWriter, r *http.Request) {
	var userDTO AddUserRequest
	if err := json.NewDecoder(r.Body).Decode(&userDTO); err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	user, err := h.bank.AddUser(r.Context(), userDTO.FullName)

	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}
}

func (h *HTTPHandlers) CreateTransaction(w http.ResponseWriter, r *http.Request) {

	idString := mux.Vars(r)["id"]

	id, err := uuid.Parse(idString)
	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	var transaction TransactionRequest
	if err := json.NewDecoder(r.Body).Decode(&transaction); err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	user, err := h.bank.CreateTransaction(r.Context(), transaction.Type, id, transaction.Amount)
	if err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		err := ErrDTO{
			error: err,
			time:  time.Now(),
		}

		http.Error(w, err.ToString(), http.StatusBadRequest)
		return
	}
}
