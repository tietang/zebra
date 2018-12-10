package main

import (
    "fmt"
)

func fibonacci(n int, start uint) []int {
    seq := make([]int, 0)
    x := n
    seq = append(seq, x)

    for {
        x = int(float64(x) * 0.618)
        if x <= int(start) {
            break
        }
        seq = append(seq, x)
    }
    return seq
}

func main() {

    n := 5000
    seq := fibonacci(n, 2)
    fmt.Println(seq)
    fmt.Println(len(seq))

    //w := &requester.Work{
    //    Request: nil,
    //    N:       1,
    //    C:       1,
    //    QPS:     100,
    //    Timeout: 3,
    //    H2:      false,
    //    //ProxyAddr:          proxyURL,
    //    Output: "csv",
    //}
    //c := make(chan os.Signal, 1)
    //signal.Notify(c, os.Interrupt)
    //go func() {
    //    <-c
    //    w.Finish()
    //    os.Exit(1)
    //}()
    //w.Run()
}
