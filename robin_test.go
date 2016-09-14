package zuul

import (
    "fmt"
    "testing"

    "github.com/tietang/assert"
    "github.com/tietang/stats"
)

func TestRobinTheSameWeight(t *testing.T) {
    var hosts1 = make([]HostInstance, 0)

    hosts1 = append(hosts1, HostInstance{Name: "a", Weight: 1})
    hosts1 = append(hosts1, HostInstance{Name: "b", Weight: 1})
    hosts1 = append(hosts1, HostInstance{Name: "c", Weight: 1})

    r := Robin{}

    h := r.Next(hosts1)
    assert.Equal(t, h.Name, hosts1[0].Name)
    h = r.Next(hosts1)
    assert.Equal(t, h.Name, hosts1[1].Name)
    h = r.Next(hosts1)
    assert.Equal(t, h.Name, hosts1[2].Name)
    h = r.Next(hosts1)
    assert.Equal(t, h.Name, hosts1[0].Name)
    times := 100
    size := len(hosts1) * times
    counter := stats.NewCounter()

    for i := 0; i < size; i++ {
        h = r.Next(hosts1)
        counter.Incr(h.Name, 1)
    }

    assert.Equal(t, int64(100), counter.Get(hosts1[0].Name).Count)
    assert.Equal(t, int64(100), counter.Get(hosts1[1].Name).Count)
    assert.Equal(t, int64(100), counter.Get(hosts1[2].Name).Count)

    fmt.Println(h.Name, h.Weight, hosts1)

    Convey()

}

func TestRobinDiffWeight(t *testing.T) {
    var hosts1 = make([]HostInstance, 0)

    hosts1 = append(hosts1, HostInstance{Name: "a", Weight: 1})
    hosts1 = append(hosts1, HostInstance{Name: "b", Weight: 2})
    hosts1 = append(hosts1, HostInstance{Name: "c", Weight: 3})

    r := Robin{}

    times := 100
    //总执行次数
    size := getTotalWeight(hosts1) * times
    counter := stats.NewCounter()
    var h *HostInstance
    for i := 0; i < size; i++ {
        h = r.Next(hosts1)
        counter.Incr(h.Name, 1)
    }
    //各个权重命中是否正确
    assert.Equal(t, int64(100 * hosts1[0].Weight), counter.Get(hosts1[0].Name).Count)
    assert.Equal(t, int64(100 * hosts1[1].Weight), counter.Get(hosts1[1].Name).Count)
    assert.Equal(t, int64(100 * hosts1[2].Weight), counter.Get(hosts1[2].Name).Count)

}

func getTotalWeight(hosts []HostInstance) int {
    c := 0
    for _, v := range hosts {
        c = c + v.Weight
    }
    return c
}
