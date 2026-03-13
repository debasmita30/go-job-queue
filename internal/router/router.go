package router

import (
	"database/sql"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/debasmita30/go-job-queue/internal/handlers"
	"github.com/debasmita30/go-job-queue/internal/queue"
)

func Setup(db *sql.DB, dispatcher *queue.Dispatcher) *gin.Engine {
	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowAllOrigins:  true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
	}))

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
```

---

**`go.mod`**
```
module github.com/debasmita30/go-job-queue

go 1.21

require (
	github.com/gin-contrib/cors v1.4.0
	github.com/gin-gonic/gin v1.9.1
	github.com/google/uuid v1.4.0
	github.com/joho/godotenv v1.5.1
	github.com/lib/pq v1.10.9
)
