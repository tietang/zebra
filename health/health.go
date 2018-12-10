package health

const (
    STATUS_UP   = "UP"
    STATUS_DOWN = "DOWN"
)

type HealthChecker interface {
    CheckHealth(rootHealth *RootHealth)
}

type Health struct {
    Status  string             `json:"status"`
    Desc    string             `json:"desc"`
    Healths map[string]*Health `json:"states,omitempty"`
}

type RootHealth struct {
    Health
    HealthCheckers []HealthChecker `json:"-"`
}

func (h *RootHealth) Check() {
    for _, hc := range h.HealthCheckers {
        hc.CheckHealth(h)
    }
}

func (h *RootHealth) Add(healthChecker HealthChecker) {
    h.HealthCheckers = append(h.HealthCheckers, healthChecker)
}

func (h *RootHealth) AddAll(healthCheckers []HealthChecker) {
    for _, healthChecker := range healthCheckers {
        h.Add(healthChecker)
    }
}
