package main

import (
	"context"
	"minibank/internal/adapters/api"
	"minibank/internal/adapters/config"
	"minibank/internal/adapters/db"
	domains "minibank/internal/domain/users"
	"minibank/internal/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
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

	storage := db.NewUserRepository(pool, log)        // db
	bankService := domains.NewBankService(storage)    // usecases
	handlers := api.NewHTTPHandlers(bankService, log) // api
	router := api.NewRouter(handlers)                 //api

	srv := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: router,
	}

	log.Info("http server starting", "addr", ":"+cfg.Port)

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("http server failed", "err", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	<-stop

	log.Info("shutting down server...")
	shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutCtx); err != nil {
		log.Error("graceful shutdown failed", "err", err)
	}

}
