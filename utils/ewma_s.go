package utils

import (
    "github.com/rcrowley/go-metrics"
    "math"
    "time"
)

func NewEwma(seconds int) metrics.EWMA {
    ewma := metrics.NewEWMA(1 - math.Exp(-5.0/60.0/(1/float64(seconds))))
    ticker := time.NewTicker(5e9)
    go func() {
        for {
            select {
            case <-ticker.C:
                ewma.Tick()
            }
        }
    }()

    return ewma
}
