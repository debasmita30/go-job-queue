package config

import (
    "os"
    "strconv"
    "github.com/joho/godotenv"
)

type Config struct {
    Port        string
    DatabaseURL string
    WorkerCount int
}

func Load() *Config {
    godotenv.Load()

    workerCount, err := strconv.Atoi(os.Getenv("WORKER_COUNT"))
    if err != nil || workerCount == 0 {
        workerCount = 5
    }

    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    return &Config{
        Port:        port,
        DatabaseURL: os.Getenv("DATABASE_URL"),
        WorkerCount: workerCount,
    }
}