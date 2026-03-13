package main

import (
	"context"
	"minibank/internal/adapters/api"
	"minibank/internal/adapters/config"
	"minibank/internal/adapters/db"
	"minibank/internal/adapters/logger"
	domains "minibank/internal/domain/user"
	"net/http"

	"github.com/joho/godotenv"
)

func main() {
	ctx := context.Background()
	godotenv.Load(".env")

	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	log := logger.NewLogger(cfg.LogLevel)

	log.Info("application started")

	pool, err := db.NewConnection(ctx, cfg.DataBaseURL)
	if err != nil {
		log.Error("failed to connect database", "err", err)
		return
	}
	defer pool.Close()

	log.Info("database connected")

	storage := db.NewUserRepository(pool, log)
	bankService := domains.NewBankService(storage)
	handlers := api.NewHTTPHandlers(bankService, log)
	router := api.NewRouter(handlers)

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Info("http server starting", "addr", ":"+cfg.Port)

	if err := srv.ListenAndServe(); err != nil {
		log.Error("http server failed", "err", err)
	}
}
