package main

import (
    "fmt"
    "github.com/rcrowley/go-metrics"
    "log"
    "math/rand"
    "os"
    "time"
)

func main() {

    s := metrics.NewExpDecaySample(1028, 0.015) // or metrics.NewUniformSample(1028)
    h := metrics.NewHistogram(s)
    metrics.Register("baz", h)
    fd, _ := os.OpenFile("a.csv", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)
    fmt.Println(fd.Name())
    i := 0
    go metrics.Log(metrics.DefaultRegistry, 5*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
    for {
        seed := uint64(rand.Uint64() % 200)
        if seed == 0 {
            seed = 1
        }
        s := time.Now().Second() % 7
        seed = uint64(10*s) + uint64(rand.Uint64()%seed)
        h.Update(int64(seed))
        time.Sleep(time.Millisecond * 1)
        i++
        if i%1000 == 0 {
            ps := h.Percentiles([]float64{0.5, 0.75, 0.95, 0.99, 0.999})
            fd.WriteString(fmt.Sprintf("%9d,", h.Count()))
            fd.WriteString(fmt.Sprintf("%9d,", h.Min()))
            fd.WriteString(fmt.Sprintf("%9d,", h.Max()))
            fd.WriteString(fmt.Sprintf("%12.2f,", h.Mean()))
            fd.WriteString(fmt.Sprintf("%12.2f,", h.StdDev()))
            fd.WriteString(fmt.Sprintf("%12.2f,", h.Variance()))
            fd.WriteString(fmt.Sprintf("%12.2f,", ps[0]))
            fd.WriteString(fmt.Sprintf(" %12.2f,", ps[1]))
            fd.WriteString(fmt.Sprintf(" %12.2f,", ps[2]))
            fd.WriteString(fmt.Sprintf(" %12.2f,", ps[3]))
            fd.WriteString(fmt.Sprintf(" %12.2f,", ps[4]))
            fd.WriteString("\n")
        }
    }
    //
    //m := router.NewMeter(5, 5)
    //
    //i := 0
    ////go metrics.LogScaled(metrics.DefaultRegistry, 1*time.Second, 5*time.Minute, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
    //for {
    //    seed := int64(rand.Uint64() % 100)
    //    m.Mark(seed)
    //    i++
    //    time.Sleep(time.Millisecond * 1)
    //
    //    if i%1000 == 0 {
    //        fmt.Println(m.Avg(), m.Total(), m.Count())
    //    }
    //
    //}

}
