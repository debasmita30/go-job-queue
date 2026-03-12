package queue

import (
    "database/sql"
    "encoding/json"
    "time"
    "github.com/debasmita30/go-job-queue/internal/models"
)

type Queue struct {
    db *sql.DB
}

func NewQueue(db *sql.DB) *Queue {
    return &Queue{db: db}
}

func (q *Queue) Enqueue(req *models.CreateJobRequest) (*models.Job, error) {
    if req.Priority == 0 {
        req.Priority = 1
    }
    if req.MaxAttempts == 0 {
        req.MaxAttempts = 3
    }

    scheduledAt := time.Now()
    if req.ScheduledAt != nil {
        scheduledAt = *req.ScheduledAt
    }

    var job models.Job
    query := `
        INSERT INTO jobs (type, payload, priority, max_attempts, scheduled_at)
        VALUES ($1, $2, $3, $4, $5)
        RETURNING id, type, payload, status, priority, attempts, max_attempts,
                  error, created_at, updated_at, scheduled_at, processed_at
    `
    row := q.db.QueryRow(query, req.Type, req.Payload, req.Priority, req.MaxAttempts, scheduledAt)
    err := scanJob(row, &job)
    return &job, err
}

func (q *Queue) Dequeue() (*models.Job, error) {
    var job models.Job
    query := `
        UPDATE jobs SET status = 'processing', updated_at = NOW(), attempts = attempts + 1
        WHERE id = (
            SELECT id FROM jobs
            WHERE status = 'pending'
            AND scheduled_at <= NOW()
            ORDER BY priority DESC, created_at ASC
            FOR UPDATE SKIP LOCKED
            LIMIT 1
        )
        RETURNING id, type, payload, status, priority, attempts, max_attempts,
                  error, created_at, updated_at, scheduled_at, processed_at
    `
    row := q.db.QueryRow(query)
    err := scanJob(row, &job)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &job, err
}

func (q *Queue) MarkCompleted(jobID string) error {
    now := time.Now()
    _, err := q.db.Exec(`
        UPDATE jobs SET status = 'completed', processed_at = $1, updated_at = NOW()
        WHERE id = $2
    `, now, jobID)
    return err
}

func (q *Queue) MarkFailed(jobID string, errMsg string) error {
    _, err := q.db.Exec(`
        UPDATE jobs
        SET status = CASE WHEN attempts >= max_attempts THEN 'dead' ELSE 'pending' END,
            error = $1,
            updated_at = NOW(),
            scheduled_at = CASE WHEN attempts >= max_attempts THEN scheduled_at
                           ELSE NOW() + (attempts * interval '30 seconds') END
        WHERE id = $2
    `, errMsg, jobID)
    return err
}

func (q *Queue) MoveToDead(job *models.Job) error {
    payload, _ := json.Marshal(job.Payload)
    _, err := q.db.Exec(`
        INSERT INTO dead_letter_jobs (original_job_id, type, payload, error, attempts)
        VALUES ($1, $2, $3, $4, $5)
    `, job.ID, job.Type, payload, job.Error, job.Attempts)
    return err
}

func (q *Queue) GetJob(jobID string) (*models.Job, error) {
    var job models.Job
    query := `
        SELECT id, type, payload, status, priority, attempts, max_attempts,
               error, created_at, updated_at, scheduled_at, processed_at
        FROM jobs WHERE id = $1
    `
    row := q.db.QueryRow(query, jobID)
    err := scanJob(row, &job)
    if err == sql.ErrNoRows {
        return nil, nil
    }
    return &job, err
}

func (q *Queue) GetStats() (*models.QueueStats, error) {
    stats := &models.QueueStats{}
    rows, err := q.db.Query(`
        SELECT status, COUNT(*) FROM jobs GROUP BY status
    `)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    for rows.Next() {
        var status string
        var count int
        rows.Scan(&status, &count)
        stats.TotalJobs += count
        switch status {
        case "pending":
            stats.PendingJobs = count
        case "processing":
            stats.ProcessingJobs = count
        case "completed":
            stats.CompletedJobs = count
        case "failed":
            stats.FailedJobs = count
        case "dead":
            stats.DeadJobs = count
        }
    }

    var deadCount int
    q.db.QueryRow(`SELECT COUNT(*) FROM dead_letter_jobs`).Scan(&deadCount)
    stats.DeadJobs += deadCount

    return stats, nil
}

func (q *Queue) ListJobs(status string, limit, offset int) ([]models.Job, error) {
    query := `
        SELECT id, type, payload, status, priority, attempts, max_attempts,
               error, created_at, updated_at, scheduled_at, processed_at
        FROM jobs
    `
    args := []interface{}{}
    if status != "" {
        query += ` WHERE status = $1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`
        args = append(args, status, limit, offset)
    } else {
        query += ` ORDER BY created_at DESC LIMIT $1 OFFSET $2`
        args = append(args, limit, offset)
    }

    rows, err := q.db.Query(query, args...)
    if err != nil {
        return nil, err
    }
    defer rows.Close()

    var jobs []models.Job
    for rows.Next() {
        var job models.Job
        scanJobRow(rows, &job)
        jobs = append(jobs, job)
    }
    return jobs, nil
}

func scanJob(row *sql.Row, job *models.Job) error {
    return row.Scan(
        &job.ID, &job.Type, &job.Payload, &job.Status,
        &job.Priority, &job.Attempts, &job.MaxAttempts,
        &job.Error, &job.CreatedAt, &job.UpdatedAt,
        &job.ScheduledAt, &job.ProcessedAt,
    )
}

func scanJobRow(rows *sql.Rows, job *models.Job) error {
    return rows.Scan(
        &job.ID, &job.Type, &job.Payload, &job.Status,
        &job.Priority, &job.Attempts, &job.MaxAttempts,
        &job.Error, &job.CreatedAt, &job.UpdatedAt,
        &job.ScheduledAt, &job.ProcessedAt,
    )
}