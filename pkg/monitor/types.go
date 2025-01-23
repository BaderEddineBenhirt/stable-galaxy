package monitor

import (
    "time"
)

type Status string

const (
    StatusHealthy Status = "healthy"
    StatusWarning Status = "warning"
    StatusError   Status = "error"
)

type Version struct {
    Number      string
    Status      Status
    Metrics     *Metrics
    LastChecked time.Time
    Errors      []Error
}

type Metrics struct {
    CPUUsage    float64
    MemoryUsage float64
    ErrorRate   float64
    Latency     time.Duration
}

type Error struct {
    Message   string
    Timestamp time.Time
    Severity  string
}

