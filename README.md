# GoQueue — Background Job Processing Engine

> **Production-grade distributed job queue built in Go** — handles background tasks so your app stays fast and responsive.

---

## 🔴 The Problem

Modern web applications need to handle slow, resource-intensive tasks like:

- Sending 10,000 transactional emails after a flash sale
- Generating PDF reports from millions of database rows
- Delivering webhooks to 50 third-party services simultaneously
- Processing payment notifications from Razorpay / Stripe

**If you process these tasks synchronously**, your API hangs, users wait, and timeouts cascade into failures.

```
❌ Without a job queue:
User clicks "Generate Report" → API blocks for 45 seconds → User sees timeout → Data lost
```

---

## ✅ The Solution

GoQueue decouples task submission from task execution using a persistent queue and concurrent worker pool.

```
✅ With GoQueue:
User clicks "Generate Report" → API returns instantly ("Job queued!") → 
5 workers process in background → User notified when done
```

**Key guarantees:**
- Jobs are never lost — stored in PostgreSQL before any worker touches them
- Auto-retry with exponential backoff if a job fails
- Dead letter queue captures jobs that exhaust all retries
- Priority scheduling — urgent jobs jump the queue
- 5 goroutines process jobs truly in parallel

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                        CLIENT / API                             │
│                   POST /api/v1/jobs                             │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                    GIN HTTP SERVER                              │
│              JobHandler · StatsHandler                          │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                     DISPATCHER                                  │
│         Enqueues jobs → PostgreSQL (FOR UPDATE SKIP LOCKED)     │
└──────┬──────────┬──────────┬──────────┬──────────┬─────────────┘
       │          │          │          │          │
       ▼          ▼          ▼          ▼          ▼
   Worker 0   Worker 1   Worker 2   Worker 3   Worker 4
  (goroutine)(goroutine)(goroutine)(goroutine)(goroutine)
       │          │          │          │          │
       └──────────┴──────────┴──────────┴──────────┘
                            │
               ┌────────────┴────────────┐
               ▼                         ▼
        ✅ completed              ❌ failed → retry
                                         │
                                  max_attempts exceeded?
                                         │
                                         ▼
                                  💀 dead letter queue
```

```
Job Lifecycle:

[pending] → [processing] → [completed]
               │
               └──(on error)──► [failed] ──(retry)──► [processing]
                                              │
                                    (attempts >= max)
                                              │
                                              ▼
                                           [dead]
```

---

## 📁 Project Structure

```
go-job-queue/
│
├── cmd/
│   └── server/
│       └── main.go                 # Entry point — starts server + workers
│
├── internal/
│   ├── config/
│   │   └── config.go               # Env config (DATABASE_URL, PORT, WORKER_COUNT)
│   │
│   ├── database/
│   │   └── database.go             # PostgreSQL connection + migrations
│   │
│   ├── models/
│   │   └── job.go                  # Job struct — status, priority, payload, error
│   │
│   ├── queue/
│   │   ├── queue.go                # Dispatcher — enqueue + fetch with SKIP LOCKED
│   │   └── worker.go               # Worker goroutines — process + retry logic
│   │
│   ├── handlers/
│   │   └── job_handler.go          # REST handlers — enqueue, list, get, stats
│   │
│   └── router/
│       └── router.go               # Gin router + CORS middleware
│
├── migrations/
│   └── 001_create_jobs.sql         # PostgreSQL schema
│
├── index.html                      # Live demo frontend
├── Dockerfile                      # Multi-stage Docker build
├── docker-compose.yml              # Local dev with Postgres
├── go.mod
├── go.sum
├── .env.example
└── .gitignore
```

---

## 🚀 API Reference

### Health Check
```http
GET /health
```
```json
{ "status": "healthy", "service": "go-job-queue", "version": "1.0.0" }
```

---

### Enqueue a Job
```http
POST /api/v1/jobs
Content-Type: application/json
```
```json
{
  "type": "send_email",
  "payload": { "to": "user@example.com", "subject": "Welcome!" },
  "priority": 3,
  "max_attempts": 3
}
```
**Response:**
```json
{
  "job": {
    "id": "550e8400-e29b-41d4-a716-446655440000",
    "type": "send_email",
    "status": "pending",
    "priority": 3,
    "attempts": 0,
    "max_attempts": 3,
    "created_at": "2026-03-13T09:43:27Z"
  }
}
```

---

### List Jobs
```http
GET /api/v1/jobs?status=pending&limit=20
```

| Query Param | Options | Default |
|---|---|---|
| `status` | `pending`, `processing`, `completed`, `failed`, `dead` | all |
| `limit` | 1–100 | 20 |

---

### Get Single Job
```http
GET /api/v1/jobs/:id
```

---

### Queue Statistics
```http
GET /api/v1/stats
```
```json
{
  "total_jobs": 142,
  "pending_jobs": 3,
  "processing_jobs": 2,
  "completed_jobs": 134,
  "failed_jobs": 1,
  "dead_jobs": 2
}
```

---

## ⚙️ Supported Job Types

| Type | Payload Fields | Description |
|---|---|---|
| `send_email` | `to`, `subject`, `body` | Background email delivery |
| `generate_report` | `report_type`, `format` | Async report generation |
| `webhook_delivery` | `url`, `event`, `data` | Reliable webhook dispatch |

---

## 🔁 Retry Logic

GoQueue uses **exponential backoff** to retry failed jobs without hammering downstream services:

```
Attempt 1 fails → wait 30s  → retry
Attempt 2 fails → wait 60s  → retry
Attempt 3 fails → wait 90s  → marked DEAD
```

Workers use `FOR UPDATE SKIP LOCKED` — a PostgreSQL feature that prevents two workers from picking the same job even under high concurrency.

---

## 🖥️ Running Locally

### With Docker Compose (recommended)
```bash
git clone https://github.com/debasmita30/go-job-queue
cd go-job-queue
cp .env.example .env
docker-compose up --build
```

### Without Docker
```bash
# Requires Go 1.23+ and PostgreSQL running locally
export DATABASE_URL=postgres://postgres:password@localhost:5432/jobqueue?sslmode=disable
export WORKER_COUNT=5
export PORT=8080
go run ./cmd/server
```

API available at `http://localhost:8080`

---

## 📸 Screenshots

### Live Dashboard
<!-- Add screenshot here -->
![Dashboard](screenshots/dashboard.png)

### Queue Stats
<!-- Add screenshot here -->
![Stats](screenshots/stats.png)

### Job Submission
<!-- Add screenshot here -->
![Job Form](screenshots/job-form.png)

### Jobs Table
<!-- Add screenshot here -->
![Jobs Table](screenshots/jobs-table.png)

---

## 🛠️ Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.23 |
| HTTP Framework | Gin |
| Database | PostgreSQL 15 |
| Concurrency | Goroutines + Channels |
| Containerization | Docker (multi-stage build) |
| Cloud Deployment | Render |
| Frontend | Vanilla React (CDN) |

---

## 💡 Real-World Parallels

This architecture mirrors what runs at scale in production:

- **Sidekiq** (Ruby) — Redis-backed job queue used by GitHub, Shopify
- **Celery** (Python) — Distributed task queue used by Instagram, Mozilla
- **BullMQ** (Node.js) — Queue used by Netlify, Linear
- **Faktory** — Language-agnostic job server by the creator of Sidekiq

GoQueue implements the same core patterns: persistent storage, concurrent workers, retry with backoff, and dead letter queues.

---

## 🔗 Live Demo

| Resource | Link |
|---|---|
| 🌐 Live Frontend | https://lovely-marigold-29675d.netlify.app |
| ⚙️ API Health | https://go-job-queue.onrender.com/health |
| 📊 Live Stats | https://go-job-queue.onrender.com/api/v1/stats |
| 💻 Source Code | https://github.com/debasmita30/go-job-queue |

> **Note:** Hosted on Render free tier — first request may take 30–60s to wake the server.

---

## 👩‍💻 Author

**Debasmita Chatterjee**  
[GitHub](https://github.com/debasmita30) · [Portfolio](https://leafy-cajeta-9270ea.netlify.app/)
