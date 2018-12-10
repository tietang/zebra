package meter

import (
    "errors"
    "sync"
    "sync/atomic"
    "time"
)

type Counter struct {
    Count      int64
    Value      int64
    TimeWindow int64
}

type RouteMeter struct {
    m1       map[int64]*Counter
    interval int
    saveSize int
    lock     *sync.RWMutex
}

func NewRouteMeter(intervalSeconds, saveSize int) *RouteMeter {
    m := &RouteMeter{
        m1:       make(map[int64]*Counter),
        interval: intervalSeconds,
        saveSize: saveSize,
        lock:     new(sync.RWMutex),
    }
    return m
}

func (m *RouteMeter) Mark(n int64) {
    key := m.key()
    key_3 := key - int64(m.interval*m.saveSize)

    m.lock.Lock()
    defer m.lock.Unlock()
    delete(m.m1, key_3)
    v, ok := m.m1[key]
    if !ok {
        v = &Counter{}
    }
    atomic.AddInt64(&v.Count, 1)
    atomic.AddInt64(&v.Value, n)
    v.TimeWindow = key

    m.m1[key] = v
}
func (m *RouteMeter) key() int64 {
    second := time.Now().Unix()
    key := int64(m.interval) * (second / int64(m.interval))
    return key
}
func (m *RouteMeter) Get(n int) *Counter {
    if n > 0 || n > m.saveSize {
        panic(errors.New("[n] must be more than or equal to [1 - save size], and must be less than or equal to 0. "))
    }
    m.lock.RLock()
    defer m.lock.RUnlock()
    key := m.key() + int64(n*m.interval)
    return m.m1[key]
}

func (m *RouteMeter) Snapshot(n int) Counter {
    c := m.Get(n)
    if c == nil {
        return Counter{}
    }
    return Counter{
        Count:      c.Count,
        TimeWindow: c.TimeWindow,
        Value:      c.Value,
    }
}

func (m *RouteMeter) Avg() int64 {
    c := m.Get(0)
    return c.TimeWindow / c.Count
}

func (m *RouteMeter) Total() int64 {
    c := m.Get(0)
    return c.TimeWindow
}

func (m *RouteMeter) Count() int64 {
    c := m.Get(0)
    return c.Count
}
