package router

import (
    "database/sql"
    "github.com/gin-gonic/gin"
    "github.com/debasmita30/go-job-queue/internal/handlers"
    "github.com/debasmita30/go-job-queue/internal/queue"
)

func Setup(db *sql.DB, dispatcher *queue.Dispatcher) *gin.Engine {
    r := gin.Default()

    jobHandler := handlers.NewJobHandler(dispatcher)

    r.GET("/health", jobHandler.HealthCheck)

    api := r.Group("/api/v1")
    {
        jobs := api.Group("/jobs")
        {
            jobs.POST("", jobHandler.EnqueueJob)
            jobs.GET("", jobHandler.ListJobs)
            jobs.GET("/:id", jobHandler.GetJob)
        }
        api.GET("/stats", jobHandler.GetStats)
    }

    return r
}