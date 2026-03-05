package main

import (
	"minibank/internal/adapters/api"
	domains "minibank/internal/domain/user"
	"net/http"
)

func main() {
	//ctx := context.Background()

	storage := domains.NewStorage()
	bankService := domains.NewBankService(storage)
	handlers := api.NewHTTPHandlers(bankService)
	router := api.NewRouter(handlers)

	srv := &http.Server{
		Addr:    ":9091",
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
