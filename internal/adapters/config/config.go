package config

import (
	"errors"
	"os"
)

type Config struct {
	Port        string
	DataBaseURL string
}

func Load() (Config, error) {
	port := os.Getenv("PORT")
	if port == "" {
		port = "9091"
	}

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		return Config{}, errors.New("DATABASE_URL is not set")
	}

	return Config{
		Port:        port,
		DataBaseURL: databaseURL,
	}, nil
}
