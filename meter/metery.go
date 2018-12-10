package meter

import (
    "encoding/json"
    "github.com/rcrowley/go-metrics"
    "math"
    "sync"
    "time"
)

func MarshalJSON(registries map[string]metrics.Registry) ([]byte, error) {
    data := make(map[string]interface{})
    for key, value := range registries {
        data[key] = MarshalJSONRegistry(value)
    }
    return json.Marshal(data)
}

// MarshalJSON returns a byte slice containing a JSON representation of all
// the metrics in the Registry.
func MarshalJSONRegistry(r metrics.Registry) (map[string]map[string]interface{}) {
    data := make(map[string]map[string]interface{})

    r.Each(func(name string, i interface{}) {
        values := make(map[string]interface{})
        switch metric := i.(type) {
        case metrics.Counter:
            values["count"] = metric.Count()
        case metrics.Gauge:
            values["value"] = metric.Value()
        case metrics.GaugeFloat64:
            values["value"] = metric.Value()
        case metrics.Healthcheck:
            values["error"] = nil
            metric.Check()
            if err := metric.Error(); nil != err {
                values["error"] = metric.Error().Error()
            }
        case metrics.Histogram:
            h := metric.Snapshot()
            ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
            values["count"] = h.Count()
            values["min"] = h.Min()
            values["max"] = h.Max()
            values["mean"] = h.Mean()
            values["stddev"] = h.StdDev()
            values["median"] = ps[0]
            values["75%"] = ps[1]
            values["95%"] = ps[2]
            values["99%"] = ps[3]
            values["99.9%"] = ps[4]
        case metrics.Meter:
            m := metric.Snapshot()
            values["count"] = m.Count()
            values["1m.rate"] = m.Rate1()
            values["5m.rate"] = m.Rate5()
            values["15m.rate"] = m.Rate15()
            values["mean.rate"] = m.RateMean()
        case MeterX:
            m := metric.Snapshot()
            values["count"] = m.Count()
            values["1x.rate"] = m.Rate1x()
            values["2x.rate"] = m.Rate2x()
            values["3x.rate"] = m.Rate3x()
            values["mean.rate"] = m.RateMean()
        case metrics.Timer:
            t := metric.Snapshot()
            ps := t.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
            values["count"] = t.Count()
            values["min"] = t.Min()
            values["max"] = t.Max()
            values["mean"] = t.Mean()
            values["stddev"] = t.StdDev()
            values["median"] = ps[0]
            values["75%"] = ps[1]
            values["95%"] = ps[2]
            values["99%"] = ps[3]
            values["99.9%"] = ps[4]
            values["1m.rate"] = t.Rate1()
            values["5m.rate"] = t.Rate5()
            values["15m.rate"] = t.Rate15()
            values["mean.rate"] = t.RateMean()
        }
        data[name] = values
    })
    return data
}

//var MeterSeconds map[string]int = make(map[string]int)

// Meters count events to produce exponentially-weighted moving average rates
// at one-, five-, and fifteen-minutes and a mean rate.
type MeterX interface {
    Count() int64
    Mark(int64)
    Rate1x() float64
    Rate2x() float64
    Rate3x() float64
    RateMean() float64
    Snapshot() MeterX
    Seconds() int
}

// GetOrRegisterMeter returns an existing MeterX or constructs and registers a
// new StandardMeterX.
func GetOrRegisterMeterX(name string, seconds int, r metrics.Registry) MeterX {
    if nil == r {
        r = DefaultRegistry
    }
    var m MeterX
    obj := r.Get(name)
    if obj != nil {
        m = obj.(MeterX)
        if m.Seconds() != seconds {
            r.Unregister(name)
            m = nil
        }

    }
    if m == nil {
        m = NewMeter(seconds)
    }
    r.Register(name, m)
    return m
}

// NewMeter constructs a new StandardMeterX and launches a goroutine.
func NewMeter(seconds int) *StandardMeterX {
    m := newStandardMeter(seconds)
    arbiter.Lock()
    defer arbiter.Unlock()
    arbiter.meters = append(arbiter.meters, m)
    if !arbiter.started {
        arbiter.started = true
        go arbiter.tick()
    }
    return m
}

// NewMeter constructs and registers a new StandardMeterX and launches a
// goroutine.
func NewRegisteredMeter(name string, seconds int, r metrics.Registry) MeterX {
    c := NewMeter(seconds)
    if nil == r {
        r = metrics.DefaultRegistry
    }

    r.Register(name, c)
    return c
}

// MeterSnapshot is a read-only copy of another MeterX.
type MeterSnapshot struct {
    count                            int64
    rate1x, rate2x, rate3x, rateMean float64
    seconds                          int
}

func (m *MeterSnapshot) Seconds() int {
    return m.seconds
}

// Count returns the count of events at the TimeWindow the snapshot was taken.
func (m *MeterSnapshot) Count() int64 { return m.count }

// Mark panics.
func (*MeterSnapshot) Mark(n int64) {
    panic("Mark called on a MeterSnapshot")
}

// Rate1x returns the one-minute moving average rate of events per second at the
// TimeWindow the snapshot was taken.
func (m *MeterSnapshot) Rate1x() float64 { return m.rate1x }

// Rate2x returns the five-minute moving average rate of events per second at
// the TimeWindow the snapshot was taken.
func (m *MeterSnapshot) Rate2x() float64 { return m.rate2x }

// Rate3x returns the fifteen-minute moving average rate of events per second
// at the TimeWindow the snapshot was taken.
func (m *MeterSnapshot) Rate3x() float64 { return m.rate3x }

// RateMean returns the meter's mean rate of events per second at the TimeWindow the
// snapshot was taken.
func (m *MeterSnapshot) RateMean() float64 { return m.rateMean }

// Snapshot returns the snapshot.
func (m *MeterSnapshot) Snapshot() MeterX { return m }

// StandardMeterX is the standard implementation of a MeterX.
type StandardMeterX struct {
    lock          sync.RWMutex
    snapshot      *MeterSnapshot
    a1x, a2x, a3x EWMAX
    startTime     time.Time
    seconds       int
}

//
//func newStandardMeter2(seconds int) *StandardMeterX {
//    a1x := metrics.NewEWMA(1 - math.Exp(-1.0/float64(seconds)))
//    a2x := metrics.NewEWMA(1 - math.Exp(-1.0/float64(seconds*2)))
//    a3x := metrics.NewEWMA(1 - math.Exp(-1.0/float64(seconds*3)))
//    return &StandardMeterX{
//        snapshot:  &MeterSnapshot{},
//        a1x:       a1x,
//        a2x:       a2x,
//        a3x:       a3x,
//        startTime: TimeWindow.Now(),
//        seconds:   seconds,
//    }
//}

func newStandardMeter(seconds int) *StandardMeterX {
    a1x := NewEWMAX(1 - math.Exp(-1.0/float64(seconds)))
    a2x := NewEWMAX(1 - math.Exp(-1.0/float64(seconds*2)))
    a3x := NewEWMAX(1 - math.Exp(-1.0/float64(seconds*3)))
    meter := &StandardMeterX{
        snapshot:  &MeterSnapshot{},
        a1x:       a1x,
        a2x:       a2x,
        a3x:       a3x,
        startTime: time.Now(),
        seconds:   seconds,
    }
    return meter
}

func (m *StandardMeterX) Seconds() int {
    return m.seconds
}

// Count returns the number of events recorded.
func (m *StandardMeterX) Count() int64 {
    m.lock.RLock()
    count := m.snapshot.count
    m.lock.RUnlock()
    return count
}

// Mark records the occurance of n events.
func (m *StandardMeterX) Mark(n int64) {
    m.lock.Lock()
    defer m.lock.Unlock()
    m.snapshot.count += n
    m.a1x.Update(n)
    m.a2x.Update(n)
    m.a3x.Update(n)
    m.updateSnapshot()

}

// Rate1x returns the one-minute moving average rate of events per second.
func (m *StandardMeterX) Rate1x() float64 {
    m.lock.RLock()
    rate1x := m.snapshot.rate1x
    m.lock.RUnlock()
    return rate1x
}

// Rate2x returns the five-minute moving average rate of events per second.
func (m *StandardMeterX) Rate2x() float64 {
    m.lock.RLock()
    rate2x := m.snapshot.rate2x
    m.lock.RUnlock()
    return rate2x
}

// Rate3x returns the fifteen-minute moving average rate of events per second.
func (m *StandardMeterX) Rate3x() float64 {
    m.lock.RLock()
    rate3x := m.snapshot.rate3x
    m.lock.RUnlock()
    return rate3x
}

// RateMean returns the meter's mean rate of events per second.
func (m *StandardMeterX) RateMean() float64 {
    m.lock.RLock()
    rateMean := m.snapshot.rateMean
    m.lock.RUnlock()
    return rateMean
}

// Snapshot returns a read-only copy of the meter.
func (m *StandardMeterX) Snapshot() MeterX {
    m.lock.RLock()
    snapshot := *m.snapshot
    m.lock.RUnlock()
    return &snapshot
}

func (m *StandardMeterX) updateSnapshot() {
    // should run with write lock held on m.lock
    snapshot := m.snapshot
    snapshot.rate1x = m.a1x.Rate()
    snapshot.rate2x = m.a2x.Rate()
    snapshot.rate3x = m.a3x.Rate()
    snapshot.rateMean = float64(snapshot.count) / time.Since(m.startTime).Seconds()
    snapshot.seconds = m.seconds
}

func (m *StandardMeterX) tick() {
    m.lock.Lock()
    defer m.lock.Unlock()
    m.a1x.Tick()
    m.a2x.Tick()
    m.a3x.Tick()
    m.updateSnapshot()
}

type meterArbiter struct {
    sync.RWMutex
    started bool
    meters  []*StandardMeterX
    ticker  *time.Ticker
}

var arbiter = meterArbiter{ticker: time.NewTicker(1e9)}

// Ticks meters on the scheduled interval
func (ma *meterArbiter) tick() {
    for {
        select {
        case <-ma.ticker.C:
            ma.tickMeters()
        }
    }
}

func (ma *meterArbiter) tickMeters() {
    ma.RLock()
    defer ma.RUnlock()
    for _, meter := range ma.meters {
        meter.tick()
    }
}
