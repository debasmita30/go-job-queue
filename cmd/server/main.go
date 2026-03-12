package main

import (
    "log"
    "github.com/debasmita30/go-job-queue/internal/config"
    "github.com/debasmita30/go-job-queue/internal/database"
    "github.com/debasmita30/go-job-queue/internal/queue"
    "github.com/debasmita30/go-job-queue/internal/router"
)

func main() {
    cfg := config.Load()

    db, err := database.Connect(cfg.DatabaseURL)
    if err != nil {
        log.Fatal("Failed to connect to database:", err)
    }
    defer db.Close()

    if err := database.Migrate(db); err != nil {
        log.Fatal("Failed to run migrations:", err)
    }

    dispatcher := queue.NewDispatcher(db, cfg.WorkerCount)
    dispatcher.Start()
    defer dispatcher.Stop()

    r := router.Setup(db, dispatcher)
    log.Printf("Server starting on port %s", cfg.Port)
    r.Run(":" + cfg.Port)
}