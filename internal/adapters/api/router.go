package api

import "github.com/gorilla/mux"

func NewRouter(h *HTTPHandlers) *mux.Router {
	router := mux.NewRouter()

	router.Path("/bank").Methods("GET").HandlerFunc(h.GetAllUsersH)
	router.Path("/bank/{id}").Methods("GET").HandlerFunc(h.GetUserByIDH)
	router.Path("/bank").Methods("POST").HandlerFunc(h.AddUserH)
	router.Path("/bank/{id}/transaction").Methods("POST").HandlerFunc(h.CreateTransaction)
	router.Path("/bank/{id}").Methods("DELETE").HandlerFunc(h.DeleteUserH)
	return router
}
