<div align="center">

# ⚙️ GoQueue — Background Job Processing Engine

### Never Make Your Users Wait Again — Background Jobs at Scale

[![Go](https://img.shields.io/badge/Go-1.23-00ADD8?style=flat-square&logo=go&logoColor=white)](https://go.dev)
[![Gin](https://img.shields.io/badge/Gin-Framework-00ADD8?style=flat-square&logo=go&logoColor=white)](https://gin-gonic.com)
[![PostgreSQL](https://img.shields.io/badge/PostgreSQL-15-336791?style=flat-square&logo=postgresql&logoColor=white)](https://postgresql.org)
[![Docker](https://img.shields.io/badge/Docker-Containerized-2496ED?style=flat-square&logo=docker&logoColor=white)](https://docker.com)
[![Live Demo](https://img.shields.io/badge/Live-Demo-00C851?style=flat-square&logo=render&logoColor=white)](https://lovely-marigold-29675d.netlify.app)

<br/>

[![Typing SVG](https://readme-typing-svg.demolab.com?font=Fira+Code&size=20&duration=3000&pause=800&color=00FF88&center=true&vCenter=true&multiline=false&width=750&lines=APIs+shouldn't+block+on+slow+tasks...;GoQueue+handles+them+in+the+background.;5+goroutines+processing+jobs+in+parallel;PostgreSQL-backed+%E2%80%94+zero+job+loss+guaranteed;Auto-retry+with+exponential+backoff;Priority+queues+%2B+dead+letter+support+%F0%9F%9A%80)](https://git.io/typing-svg)

<br/>

🌐 **[Live Demo](https://go-job-queue.vercel.app/)** &nbsp;|&nbsp; ⚙️ **[API Health](https://go-job-queue.onrender.com/health)** &nbsp;|&nbsp; 📊 **[Live Stats](https://go-job-queue.onrender.com/api/v1/stats)** &nbsp;|&nbsp; 💻 **[Source Code](https://github.com/debasmita30/go-job-queue)**

</div>

---

## 🧠 What Is This?

Imagine you run an e-commerce platform. When a customer places an order, your app needs to:

- Send a confirmation email
- Update inventory levels
- Ping the warehouse system
- Deliver a webhook to your CRM
- Generate a PDF receipt

**If you do all of this synchronously**, the customer waits 10+ seconds. If anything fails, the whole request fails.

**GoQueue solves this.** It accepts the task instantly, returns a response to the user in milliseconds, and processes everything in the background using concurrent Go goroutines — with automatic retries if anything goes wrong.

---

## 🎯 The Problem vs Solution

| Without a Job Queue | With GoQueue |
|---|---|
| ⏳ API blocks until slow task completes | ⚡ Returns instantly — job queued in < 1ms |
| 💥 One failure crashes the whole request | 🔁 Auto-retry with exponential backoff |
| 🗑️ Failed tasks silently disappear | 💀 Dead letter queue captures every failure |
| 🐌 Tasks processed one at a time | ⚡ 5 goroutines run jobs in true parallel |
| 📭 No visibility into task status | 📊 Real-time stats + job lifecycle tracking |
| 🔄 No prioritization | 🎯 Priority 1–3 — urgent jobs jump the queue |

---

## ✨ How It Works

```
Client hits POST /api/v1/jobs  →  Job saved to PostgreSQL instantly
→  Worker picks it up  →  Processes in goroutine
→  Succeeds: marked complete  |  Fails: retried with backoff
→  Max retries hit: moved to dead letter queue
```

**6-step job lifecycle — fully automatic, zero manual intervention.**

---

## 🏗️ Architecture

```
┌─────────────────────────────────────────────────────────────────┐
│                     CLIENT / FRONTEND                           │
│              POST /api/v1/jobs  ·  GET /api/v1/stats            │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                  GIN HTTP SERVER + CORS                         │
│         JobHandler · StatsHandler · HealthCheck                 │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                      DISPATCHER                                 │
│   Persists job to PostgreSQL  ·  FOR UPDATE SKIP LOCKED         │
└──────┬──────────┬──────────┬──────────┬──────────┬─────────────┘
       │          │          │          │          │
       ▼          ▼          ▼          ▼          ▼
   Worker 0   Worker 1   Worker 2   Worker 3   Worker 4
  (goroutine)(goroutine)(goroutine)(goroutine)(goroutine)
       │
       ├── ✅ Success → status: completed
       │
       └── ❌ Failure → exponential backoff retry
                              │
                    attempts >= max_attempts?
                              │
                              ▼
                       💀 status: dead
```

### Job Lifecycle

```
[pending] ──► [processing] ──► [completed]
                   │
                   └──(error)──► [failed] ──(retry)──► [processing]
                                               │
                                     (max attempts hit)
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
│       └── main.go                  # Entry point — starts HTTP server + worker pool
│
├── internal/
│   ├── config/
│   │   └── config.go                # Env config: DATABASE_URL, PORT, WORKER_COUNT
│   │
│   ├── database/
│   │   └── database.go              # PostgreSQL connection + schema migration
│   │
│   ├── models/
│   │   └── job.go                   # Job struct: id, type, status, payload, priority, error
│   │
│   ├── queue/
│   │   ├── queue.go                 # Dispatcher — enqueue + fetch (SKIP LOCKED)
│   │   └── worker.go                # Worker goroutines — process + retry logic
│   │
│   ├── handlers/
│   │   └── job_handler.go           # REST handlers: enqueue, list, get, stats
│   │
│   └── router/
│       └── router.go                # Gin router setup + CORS middleware
│
├── migrations/
│   └── 001_create_jobs.sql          # PostgreSQL jobs table schema
│
├── index.html                       # Live demo frontend (React CDN)
├── Dockerfile                       # Multi-stage Docker build
├── docker-compose.yml               # Local dev: app + PostgreSQL
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
  "payload": { "to": "user@example.com", "subject": "Order Confirmed!" },
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

| Param | Options | Default |
|---|---|---|
| `status` | `pending` `processing` `completed` `failed` `dead` | all |
| `limit` | 1–100 | 20 |

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

| Type | Payload | Description |
|---|---|---|
| `send_email` | `to`, `subject`, `body` | Background email delivery |
| `generate_report` | `report_type`, `format` | Async report generation |
| `webhook_delivery` | `url`, `event`, `data` | Reliable webhook dispatch |

---

## 🔁 Retry Logic

GoQueue uses **exponential backoff** — failed jobs wait longer between each retry to avoid hammering downstream services:

```
Attempt 1 fails → wait 30s  → retry
Attempt 2 fails → wait 60s  → retry  
Attempt 3 fails → wait 90s  → marked DEAD (moved to dead letter queue)
```

Workers use PostgreSQL's `FOR UPDATE SKIP LOCKED` — a battle-tested pattern that prevents two workers from ever picking the same job, even at high concurrency.

---

## 📸 Screenshots

> **Live Dashboard — Server Online**

<!-- Add screenshot here -->
&nbsp;

> **Real-time Queue Stats**

<!-- Add screenshot here -->
&nbsp;

> **Job Submission Form + Demo Scenarios**

<!-- Add screenshot here -->
&nbsp;

> **Jobs Table with Status Badges**

<!-- Add screenshot here -->
&nbsp;

---

## 🛠️ Tech Stack

| Layer | Technology | Purpose |
|---|---|---|
| Language | Go 1.23 | Performance + goroutines |
| HTTP Framework | Gin | REST API + routing |
| Database | PostgreSQL 15 | Persistent job storage |
| Concurrency | Goroutines + Channels | Parallel worker pool |
| Containerization | Docker (multi-stage) | Reproducible builds |
| Cloud | Render | Production deployment |
| Frontend | React (CDN) | Live demo dashboard |

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
# Requires Go 1.23+ and a running PostgreSQL instance
export DATABASE_URL=postgres://postgres:password@localhost:5432/jobqueue?sslmode=disable
export WORKER_COUNT=5
export PORT=8080
go run ./cmd/server
```

API live at `http://localhost:8080` · Try `GET /health` to verify.

---

## 💡 Real-World Parallels

GoQueue implements the same architecture patterns used by job queues at scale in production:

| Library | Language | Used By |
|---|---|---|
| **Sidekiq** | Ruby | GitHub, Shopify, Gitlab |
| **Celery** | Python | Instagram, Mozilla, Zapier |
| **BullMQ** | Node.js | Netlify, Linear |
| **Faktory** | Language-agnostic | Creator of Sidekiq |

Core patterns implemented: persistent job storage, concurrent worker pool, retry with exponential backoff, dead letter queue, priority scheduling.

---

## 📈 Performance

| Metric | Value |
|---|---|
| Job enqueue latency | < 5ms |
| Worker poll interval | 1 second |
| Max concurrent jobs | 5 (configurable) |
| Retry backoff formula | `30s × attempt_number` |
| Job storage | Persistent (survives restarts) |

---

## 🗺️ Roadmap

- [x] PostgreSQL-backed persistent job queue
- [x] 5 concurrent goroutine workers
- [x] Priority scheduling (1–3)
- [x] Exponential backoff retry
- [x] Dead letter queue
- [x] REST API (enqueue, list, get, stats)
- [x] CORS support
- [x] Docker + docker-compose
- [x] Deployed on Render
- [x] Live React dashboard
- [ ] Redis backend option
- [ ] Job scheduling (cron-style)
- [ ] Webhook callbacks on completion
- [ ] Prometheus metrics endpoint

---

## 👩‍💻 Author

<div align="center">

**Debasmita Chatterjee**

Backend Engineer · Go · Python · ML Systems · Cloud Deployment

[![LinkedIn](https://img.shields.io/badge/LinkedIn-Connect-0077B5?style=flat-square&logo=linkedin)](https://www.linkedin.com/in/debasmita-chatterjee/)
[![GitHub](https://img.shields.io/badge/GitHub-Follow-181717?style=flat-square&logo=github)](https://github.com/debasmita30)
[![Portfolio](https://img.shields.io/badge/Portfolio-Visit-00FF88?style=flat-square)](https://leafy-cajeta-9270ea.netlify.app/)

</div>

---

## 📄 License

This project is licensed under the MIT License.

---

<div align="center">

⭐ **If this project helped you, give it a star!**

Built with ⚙️ Go · 🐘 PostgreSQL · 🐳 Docker · ☁️ Render

</div>
