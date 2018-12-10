package meter

import (
    "github.com/rcrowley/go-metrics"
    "reflect"
    "strings"
    "sync"
)

// The standard implementation of a metrics.Registry is a mutex-protected map
// of names to metrics.
type StandardRegistry struct {
    metrics map[string]interface{}
    mutex   sync.Mutex
}

// Create a new registry.
func NewRegistry() metrics.Registry {
    return &StandardRegistry{metrics: make(map[string]interface{})}
}

// Call the given function for each registered metric.
func (r *StandardRegistry) Each(f func(string, interface{})) {
    for name, i := range r.registered() {
        f(name, i)
    }
}

// Get the metric by the given name or nil if none is registered.
func (r *StandardRegistry) Get(name string) interface{} {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    return r.metrics[name]
}

// GetAll metrics in the Registry
func (r *StandardRegistry) GetAll() map[string]map[string]interface{} {
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

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *StandardRegistry) GetOrRegister(name string, i interface{}) interface{} {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    if metric, ok := r.metrics[name]; ok {
        return metric
    }
    if v := reflect.ValueOf(i); v.Kind() == reflect.Func {
        i = v.Call(nil)[0].Interface()
    }
    r.register(name, i)
    return i
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func (r *StandardRegistry) Register(name string, i interface{}) error {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    return r.register(name, i)
}

// Run all registered healthchecks.
func (r *StandardRegistry) RunHealthchecks() {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    for _, i := range r.metrics {
        if h, ok := i.(metrics.Healthcheck); ok {
            h.Check()
        }
    }
}

// Unregister the metric with the given name.
func (r *StandardRegistry) Unregister(name string) {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    delete(r.metrics, name)
}

// Unregister all metrics.  (Mostly for testing.)
func (r *StandardRegistry) UnregisterAll() {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    for name, _ := range r.metrics {
        delete(r.metrics, name)
    }
}

func (r *StandardRegistry) register(name string, i interface{}) error {
    if _, ok := r.metrics[name]; ok {
        return metrics.DuplicateMetric(name)
    }
    switch i.(type) {
    case MeterX, metrics.Counter, metrics.Gauge, metrics.GaugeFloat64, metrics.Healthcheck, metrics.Histogram, metrics.Meter, metrics.Timer:
        r.metrics[name] = i
    }
    return nil
}

func (r *StandardRegistry) registered() map[string]interface{} {
    r.mutex.Lock()
    defer r.mutex.Unlock()
    metrics := make(map[string]interface{}, len(r.metrics))
    for name, i := range r.metrics {
        metrics[name] = i
    }
    return metrics
}

type PrefixedRegistry struct {
    underlying metrics.Registry
    prefix     string
}

func NewPrefixedRegistry(prefix string) metrics.Registry {
    return &PrefixedRegistry{
        underlying: NewRegistry(),
        prefix:     prefix,
    }
}

func NewPrefixedChildRegistry(parent metrics.Registry, prefix string) metrics.Registry {
    return &PrefixedRegistry{
        underlying: parent,
        prefix:     prefix,
    }
}

// Call the given function for each registered metric.
func (r *PrefixedRegistry) Each(fn func(string, interface{})) {
    wrappedFn := func(prefix string) func(string, interface{}) {
        return func(name string, iface interface{}) {
            if strings.HasPrefix(name, prefix) {
                fn(name, iface)
            } else {
                return
            }
        }
    }

    baseRegistry, prefix := findPrefix(r, "")
    baseRegistry.Each(wrappedFn(prefix))
}

func findPrefix(registry metrics.Registry, prefix string) (metrics.Registry, string) {
    switch r := registry.(type) {
    case *PrefixedRegistry:
        return findPrefix(r.underlying, r.prefix+prefix)
    case *StandardRegistry:
        return r, prefix
    }
    return nil, ""
}

// Get the metric by the given name or nil if none is registered.
func (r *PrefixedRegistry) Get(name string) interface{} {
    realName := r.prefix + name
    return r.underlying.Get(realName)
}

// GetAll metrics in the Registry
func (r *PrefixedRegistry) GetAll() map[string]map[string]interface{} {
    return r.underlying.GetAll()
}

// Gets an existing metric or registers the given one.
// The interface can be the metric to register if not found in registry,
// or a function returning the metric for lazy instantiation.
func (r *PrefixedRegistry) GetOrRegister(name string, metric interface{}) interface{} {
    realName := r.prefix + name
    return r.underlying.GetOrRegister(realName, metric)
}

// Register the given metric under the given name. The name will be prefixed.
func (r *PrefixedRegistry) Register(name string, metric interface{}) error {
    realName := r.prefix + name
    return r.underlying.Register(realName, metric)
}

// Run all registered healthchecks.
func (r *PrefixedRegistry) RunHealthchecks() {
    r.underlying.RunHealthchecks()
}

// Unregister the metric with the given name. The name will be prefixed.
func (r *PrefixedRegistry) Unregister(name string) {
    realName := r.prefix + name
    r.underlying.Unregister(realName)
}

// Unregister all metrics.  (Mostly for testing.)
func (r *PrefixedRegistry) UnregisterAll() {
    r.underlying.UnregisterAll()
}

var DefaultRegistry metrics.Registry = NewRegistry()

// Call the given function for each registered metric.
func Each(f func(string, interface{})) {
    DefaultRegistry.Each(f)
}

// Get the metric by the given name or nil if none is registered.
func Get(name string) interface{} {
    return DefaultRegistry.Get(name)
}

// Gets an existing metric or creates and registers a new one. Threadsafe
// alternative to calling Get and Register on failure.
func GetOrRegister(name string, i interface{}) interface{} {
    return DefaultRegistry.GetOrRegister(name, i)
}

// Register the given metric under the given name.  Returns a DuplicateMetric
// if a metric by the given name is already registered.
func Register(name string, i interface{}) error {
    return DefaultRegistry.Register(name, i)
}

// Register the given metric under the given name.  Panics if a metric by the
// given name is already registered.
func MustRegister(name string, i interface{}) {
    if err := Register(name, i); err != nil {
        panic(err)
    }
}

// Run all registered healthchecks.
func RunHealthchecks() {
    DefaultRegistry.RunHealthchecks()
}

// Unregister the metric with the given name.
func Unregister(name string) {
    DefaultRegistry.Unregister(name)
}
