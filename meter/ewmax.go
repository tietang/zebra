package meter

import (
    "math"
    "sync"
    "sync/atomic"
)

var UseNilMetrics bool = false

// EWMAXs continuously calculate an exponentially-weighted moving average
// based on an outside source of clock ticks.
type EWMAX interface {
    Rate() float64
    Snapshot() EWMAX
    Tick()
    Update(int64)
}

// NewEWMAX constructs a new EWMAX with the given alpha.
func NewEWMAX(alpha float64) EWMAX {
    if UseNilMetrics {
        return NilEWMAX{}
    }
    return &StandardEWMAX{alpha: alpha}
}

// NewEWMAX1 constructs a new EWMAX for a one-minute moving average.
func NewEWMAX1() EWMAX {
    return NewEWMAX(1 - math.Exp(-5.0/60.0/1))
}

// NewEWMAX5 constructs a new EWMAX for a five-minute moving average.
func NewEWMAX5() EWMAX {
    return NewEWMAX(1 - math.Exp(-5.0/60.0/5))
}

// NewEWMAX15 constructs a new EWMAX for a fifteen-minute moving average.
func NewEWMAX15() EWMAX {
    return NewEWMAX(1 - math.Exp(-5.0/60.0/15))
}

// EWMAXSnapshot is a read-only copy of another EWMAX.
type EWMAXSnapshot float64

// Rate returns the rate of events per second at the TimeWindow the snapshot was
// taken.
func (a EWMAXSnapshot) Rate() float64 { return float64(a) }

// Snapshot returns the snapshot.
func (a EWMAXSnapshot) Snapshot() EWMAX { return a }

// Tick panics.
func (EWMAXSnapshot) Tick() {
    panic("Tick called on an EWMAXSnapshot")
}

// Update panics.
func (EWMAXSnapshot) Update(int64) {
    panic("Update called on an EWMAXSnapshot")
}

// NilEWMAX is a no-op EWMAX.
type NilEWMAX struct{}

// Rate is a no-op.
func (NilEWMAX) Rate() float64 { return 0.0 }

// Snapshot is a no-op.
func (NilEWMAX) Snapshot() EWMAX { return NilEWMAX{} }

// Tick is a no-op.
func (NilEWMAX) Tick() {}

// Update is a no-op.
func (NilEWMAX) Update(n int64) {}

// StandardEWMAX is the standard implementation of an EWMAX and tracks the number
// of uncounted events and processes them on each tick.  It uses the
// sync/atomic package to manage uncounted events.
type StandardEWMAX struct {
    uncounted int64 // /!\ this should be the first member to ensure 64-bit alignment
    alpha     float64
    rate      float64
    init      bool
    mutex     sync.Mutex
}

// Rate returns the moving average rate of events per second.
func (a *StandardEWMAX) Rate() float64 {
    a.mutex.Lock()
    defer a.mutex.Unlock()
    return a.rate * float64(1e9)
}

// Snapshot returns a read-only copy of the EWMAX.
func (a *StandardEWMAX) Snapshot() EWMAX {
    return EWMAXSnapshot(a.Rate())
}

// Tick ticks the clock to update the moving average.  It assumes it is called
// every five seconds.
func (a *StandardEWMAX) Tick() {
    count := atomic.LoadInt64(&a.uncounted)
    atomic.AddInt64(&a.uncounted, -count)
    instantRate := float64(count) / float64(1e9)
    a.mutex.Lock()
    defer a.mutex.Unlock()
    if a.init {
        a.rate += a.alpha * (instantRate - a.rate)
    } else {
        a.init = true
        a.rate = instantRate
    }
}

// Update adds n uncounted events.
func (a *StandardEWMAX) Update(n int64) {
    atomic.AddInt64(&a.uncounted, n)
}
