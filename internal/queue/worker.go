package queue

import (
    "database/sql"
    "fmt"
    "log"
    "sync"
    "time"
    "encoding/json"
)

type JobHandler func(payload json.RawMessage) error

type Dispatcher struct {
    queue       *Queue
    workerCount int
    handlers    map[string]JobHandler
    jobChan     chan struct{}
    quit        chan struct{}
    wg          sync.WaitGroup
    mu          sync.RWMutex
}

func NewDispatcher(db *sql.DB, workerCount int) *Dispatcher {
    d := &Dispatcher{
        queue:       NewQueue(db),
        workerCount: workerCount,
        handlers:    make(map[string]JobHandler),
        jobChan:     make(chan struct{}, workerCount),
        quit:        make(chan struct{}),
    }
    d.registerDefaultHandlers()
    return d
}

func (d *Dispatcher) RegisterHandler(jobType string, handler JobHandler) {
    d.mu.Lock()
    defer d.mu.Unlock()
    d.handlers[jobType] = handler
}

func (d *Dispatcher) registerDefaultHandlers() {
    d.RegisterHandler("send_email", func(payload json.RawMessage) error {
        var data map[string]string
        json.Unmarshal(payload, &data)
        log.Printf("Sending email to %s with subject: %s", data["to"], data["subject"])
        time.Sleep(100 * time.Millisecond)
        return nil
    })

    d.RegisterHandler("generate_report", func(payload json.RawMessage) error {
        var data map[string]string
        json.Unmarshal(payload, &data)
        log.Printf("Generating report: %s", data["report_type"])
        time.Sleep(200 * time.Millisecond)
        return nil
    })

    d.RegisterHandler("webhook_delivery", func(payload json.RawMessage) error {
        var data map[string]string
        json.Unmarshal(payload, &data)
        log.Printf("Delivering webhook to: %s", data["url"])
        time.Sleep(150 * time.Millisecond)
        return nil
    })
}

func (d *Dispatcher) Start() {
    for i := 0; i < d.workerCount; i++ {
        d.wg.Add(1)
        go d.worker(i)
    }
    d.wg.Add(1)
    go d.poller()
    log.Printf("Dispatcher started with %d workers", d.workerCount)
}

func (d *Dispatcher) Stop() {
    close(d.quit)
    d.wg.Wait()
    log.Println("Dispatcher stopped")
}

func (d *Dispatcher) poller() {
    defer d.wg.Done()
    ticker := time.NewTicker(2 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case <-ticker.C:
            d.jobChan <- struct{}{}
        case <-d.quit:
            return
        }
    }
}

func (d *Dispatcher) worker(id int) {
    defer d.wg.Done()
    log.Printf("Worker %d started", id)

    for {
        select {
        case <-d.jobChan:
            d.processNextJob(id)
        case <-d.quit:
            log.Printf("Worker %d stopped", id)
            return
        }
    }
}

func (d *Dispatcher) processNextJob(workerID int) {
    job, err := d.queue.Dequeue()
    if err != nil {
        log.Printf("Worker %d: error dequeuing: %v", workerID, err)
        return
    }
    if job == nil {
        return
    }

    log.Printf("Worker %d: processing job %s (type: %s, attempt: %d)", workerID, job.ID, job.Type, job.Attempts)

    d.mu.RLock()
    handler, exists := d.handlers[job.Type]
    d.mu.RUnlock()

    if !exists {
        errMsg := fmt.Sprintf("no handler registered for job type: %s", job.Type)
        log.Printf("Worker %d: %s", workerID, errMsg)
        d.queue.MarkFailed(job.ID, errMsg)
        return
    }

    if err := handler(job.Payload); err != nil {
        log.Printf("Worker %d: job %s failed: %v", workerID, job.ID, err)
        job.Attempts++
        errStr := err.Error()
        job.Error = &errStr
        if job.Attempts >= job.MaxAttempts {
            d.queue.MoveToDead(job)
        }
        d.queue.MarkFailed(job.ID, err.Error())
        return
    }

    d.queue.MarkCompleted(job.ID)
    log.Printf("Worker %d: job %s completed", workerID, job.ID)
}

func (d *Dispatcher) GetQueue() *Queue {
    return d.queue
}