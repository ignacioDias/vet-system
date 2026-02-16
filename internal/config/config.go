package config

import (
	"os"
	"sync"
)

type Config struct {
	DatabaseURL string
	Environment string
	mu          sync.RWMutex
}

func LoadConfig() *Config {

	return &Config{
		DatabaseURL: os.Getenv("DATABASE_URL"),
		Environment: os.Getenv("ENV"),
	}
}
