package main

import (
	"context"
	"minibank/internal/adapters/api"
	"minibank/internal/adapters/config"
	"minibank/internal/adapters/db"
	domains "minibank/internal/domain/user"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load("../../.env")
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}

	pool, err := db.NewConnection(ctx, cfg.DataBaseURL)
	if err != nil {
		panic(err)
	}
	defer pool.Close()

	storage := db.NewUserRepository(pool)
	bankService := domains.NewBankService(storage)
	handlers := api.NewHTTPHandlers(bankService)
	router := api.NewRouter(handlers)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	if err := srv.ListenAndServe(); err != nil {
		panic(err)
	}
}
