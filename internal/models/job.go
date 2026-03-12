package models

import (
    "time"
    "encoding/json"
)

type JobStatus string

const (
    StatusPending    JobStatus = "pending"
    StatusProcessing JobStatus = "processing"
    StatusCompleted  JobStatus = "completed"
    StatusFailed     JobStatus = "failed"
    StatusDead       JobStatus = "dead"
)

type Job struct {
    ID          string          `json:"id"`
    Type        string          `json:"type"`
    Payload     json.RawMessage `json:"payload"`
    Status      JobStatus       `json:"status"`
    Priority    int             `json:"priority"`
    Attempts    int             `json:"attempts"`
    MaxAttempts int             `json:"max_attempts"`
    Error       *string          `json:"error,omitempty"`
    CreatedAt   time.Time       `json:"created_at"`
    UpdatedAt   time.Time       `json:"updated_at"`
    ScheduledAt time.Time       `json:"scheduled_at"`
    ProcessedAt *time.Time      `json:"processed_at,omitempty"`
}

type CreateJobRequest struct {
    Type        string          `json:"type" binding:"required"`
    Payload     json.RawMessage `json:"payload" binding:"required"`
    Priority    int             `json:"priority"`
    MaxAttempts int             `json:"max_attempts"`
    ScheduledAt *time.Time      `json:"scheduled_at"`
}

type QueueStats struct {
    TotalJobs      int `json:"total_jobs"`
    PendingJobs    int `json:"pending_jobs"`
    ProcessingJobs int `json:"processing_jobs"`
    CompletedJobs  int `json:"completed_jobs"`
    FailedJobs     int `json:"failed_jobs"`
    DeadJobs       int `json:"dead_jobs"`
}