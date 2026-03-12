package database

import (
    "database/sql"
    _ "github.com/lib/pq"
)

func Connect(databaseURL string) (*sql.DB, error) {
    db, err := sql.Open("postgres", databaseURL)
    if err != nil {
        return nil, err
    }
    if err := db.Ping(); err != nil {
        return nil, err
    }
    db.SetMaxOpenConns(25)
    db.SetMaxIdleConns(5)
    return db, nil
}

func Migrate(db *sql.DB) error {
    query := `
    CREATE TABLE IF NOT EXISTS jobs (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        type VARCHAR(100) NOT NULL,
        payload JSONB NOT NULL,
        status VARCHAR(20) NOT NULL DEFAULT 'pending',
        priority INTEGER NOT NULL DEFAULT 1,
        attempts INTEGER NOT NULL DEFAULT 0,
        max_attempts INTEGER NOT NULL DEFAULT 3,
        error TEXT,
        created_at TIMESTAMP NOT NULL DEFAULT NOW(),
        updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
        scheduled_at TIMESTAMP NOT NULL DEFAULT NOW(),
        processed_at TIMESTAMP
    );

    CREATE TABLE IF NOT EXISTS dead_letter_jobs (
        id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
        original_job_id UUID NOT NULL,
        type VARCHAR(100) NOT NULL,
        payload JSONB NOT NULL,
        error TEXT NOT NULL,
        attempts INTEGER NOT NULL,
        created_at TIMESTAMP NOT NULL DEFAULT NOW()
    );

    CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs(status);
    CREATE INDEX IF NOT EXISTS idx_jobs_priority ON jobs(priority DESC);
    CREATE INDEX IF NOT EXISTS idx_jobs_scheduled_at ON jobs(scheduled_at);
    `
    _, err := db.Exec(query)
    return err
}