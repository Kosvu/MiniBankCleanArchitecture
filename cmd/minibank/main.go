package main

import (
	"context"
	"minibank/internal/adapters/api"
	"minibank/internal/adapters/db"
	domains "minibank/internal/domain/user"
	"net/http"
)

func main() {
	ctx := context.Background()

	pool, err := db.NewConnection(ctx)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	storage := db.NewUserRepository(pool)
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
