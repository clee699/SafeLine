package metrics

import (
    "sync"
    "time"
)

// UsageStats tracks basic traffic and strategy hit metrics.
type UsageStats struct {
    sync.RWMutex
    TotalRequests     int
    TotalBlocked      int
    StrategyHits      map[string]int
    LastReset         time.Time
}

func NewUsageStats() *UsageStats {
    return &UsageStats{
        StrategyHits: make(map[string]int),
        LastReset:    time.Now(),
    }
}

func (u *UsageStats) RecordRequest(strategy string, blocked bool) {
    u.Lock()
    defer u.Unlock()
    u.TotalRequests++
    if blocked {
        u.TotalBlocked++
    }
    if strategy != "" {
        u.StrategyHits[strategy]++
    }
}

func (u *UsageStats) Reset() {
    u.Lock()
    u.TotalRequests = 0
    u.TotalBlocked = 0
    u.StrategyHits = make(map[string]int)
    u.LastReset = time.Now()
    u.Unlock()
}

func (u *UsageStats) Report() map[string]interface{} {
    u.RLock()
    defer u.RUnlock()
    return map[string]interface{}{
        "total_requests": u.TotalRequests,
        "total_blocked":  u.TotalBlocked,
        "strategy_hits":  u.StrategyHits,
        "last_reset":     u.LastReset,
    }
}