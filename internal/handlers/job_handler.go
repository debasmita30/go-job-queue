package handlers

import (
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    "github.com/debasmita30/go-job-queue/internal/models"
    "github.com/debasmita30/go-job-queue/internal/queue"
)

type JobHandler struct {
    dispatcher *queue.Dispatcher
}

func NewJobHandler(dispatcher *queue.Dispatcher) *JobHandler {
    return &JobHandler{dispatcher: dispatcher}
}

func (h *JobHandler) EnqueueJob(c *gin.Context) {
    var req models.CreateJobRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    job, err := h.dispatcher.GetQueue().Enqueue(&req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusCreated, gin.H{
        "message": "Job enqueued successfully",
        "job":     job,
    })
}

func (h *JobHandler) GetJob(c *gin.Context) {
    jobID := c.Param("id")
    job, err := h.dispatcher.GetQueue().GetJob(jobID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if job == nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "job not found"})
        return
    }
    c.JSON(http.StatusOK, job)
}

func (h *JobHandler) ListJobs(c *gin.Context) {
    status := c.Query("status")
    limit, _ := strconv.Atoi(c.DefaultQuery("limit", "20"))
    offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

    jobs, err := h.dispatcher.GetQueue().ListJobs(status, limit, offset)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{
        "jobs":   jobs,
        "limit":  limit,
        "offset": offset,
    })
}

func (h *JobHandler) GetStats(c *gin.Context) {
    stats, err := h.dispatcher.GetQueue().GetStats()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, stats)
}

func (h *JobHandler) HealthCheck(c *gin.Context) {
    c.JSON(http.StatusOK, gin.H{
        "status":  "healthy",
        "service": "go-job-queue",
        "version": "1.0.0",
    })
}