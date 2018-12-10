package router

import (
    "fmt"
    "github.com/tietang/stats"
    "sync"
    "testing"
)

func TestRobinTheSameWeight(t *testing.T) {
    Convey("给定同样的权重", t, func() {

        var hosts1 = make([]*HostInstance, 0)

        hosts1 = append(hosts1, &HostInstance{Name: "a", Weight: 1, lock: new(sync.Mutex)})
        hosts1 = append(hosts1, &HostInstance{Name: "b", Weight: 1, lock: new(sync.Mutex)})
        hosts1 = append(hosts1, &HostInstance{Name: "c", Weight: 1, lock: new(sync.Mutex)})

        r := WeightRobinRound{Lock: new(sync.Mutex)}
        var h *HostInstance
        Convey("next方法调用，分布准确", func() {
            h = r.Next("", hosts1)
            So(hosts1[0].Name, ShouldEqual, h.Name)
            h = r.Next("", hosts1)
            So(hosts1[1].Name, ShouldEqual, h.Name)
            h = r.Next("", hosts1)
            So(hosts1[2].Name, ShouldEqual, h.Name)
            h = r.Next("", hosts1)
            So(hosts1[0].Name, ShouldEqual, h.Name)
        })

        times := 100
        size := len(hosts1) * times
        counter := stats.NewCounter()

        for i := 0; i < size; i++ {
            h = r.Next("", hosts1)
            counter.Incr(h.Name, 1)
        }

    })

}

func TestRobinTheSameWeight2(t *testing.T) {
    Convey("给定同样的权重2", t, func() {

        var hosts1 = make([]*HostInstance, 0)

        hosts1 = append(hosts1, &HostInstance{Name: "172.16.1.248:7912", Weight: 1, lock: new(sync.Mutex)})
        hosts1 = append(hosts1, &HostInstance{Name: "172.16.1.248:7913", Weight: 1, lock: new(sync.Mutex)})
        hosts1 = append(hosts1, &HostInstance{Name: "172.16.1.248:7914", Weight: 1, lock: new(sync.Mutex)})
        r := WeightRobinRound{Lock: new(sync.Mutex)}
        var h *HostInstance
        Convey("next方法调用，分布准确2", func() {
            h = r.Next("", hosts1)
            fmt.Println(h)
            So(hosts1[0].Name, ShouldEqual, h.Name)
            h = r.Next("", hosts1)
            fmt.Println(h)
            So(hosts1[1].Name, ShouldEqual, h.Name)
            h = r.Next("", hosts1)
            fmt.Println(h)
            So(hosts1[2].Name, ShouldEqual, h.Name)
            h = r.Next("", hosts1)
            fmt.Println(h)
            So(hosts1[0].Name, ShouldEqual, h.Name)
        })

        times := 100
        size := len(hosts1) * times
        counter := stats.NewCounter()

        for i := 0; i < size; i++ {
            h = r.Next("", hosts1)
            counter.Incr(h.Name, 1)
        }
        Convey("调用100次后的统计2", func() {
            So(int64(100), ShouldEqual, counter.Get(hosts1[0].Name).Count)
            So(int64(100), ShouldEqual, counter.Get(hosts1[1].Name).Count)
            So(int64(100), ShouldEqual, counter.Get(hosts1[2].Name).Count)

            fmt.Println(h.Name, h.Weight, hosts1)
        })

    })

}

func TestRobinDiffWeight(t *testing.T) {
    Convey("权重不同", t, func() {

        var hosts1 = make([]*HostInstance, 0)

        hosts1 = append(hosts1, &HostInstance{Name: "a", Weight: 1, lock: new(sync.Mutex)})
        hosts1 = append(hosts1, &HostInstance{Name: "b", Weight: 2, lock: new(sync.Mutex)})
        hosts1 = append(hosts1, &HostInstance{Name: "c", Weight: 3, lock: new(sync.Mutex)})

        r := WeightRobinRound{Lock: new(sync.Mutex)}

        times := int32(100)
        //总执行次数
        size := getTotalWeight(hosts1) * times
        counter := stats.NewCounter()
        var h *HostInstance
        for i := 0; i < int(size); i++ {
            h = r.Next("", hosts1)
            //fmt.Println(h)
            counter.Incr(h.Name, 1)
        }
        //d, _ := json.Marshal(counter.GetAll())
        //fmt.Println(size, string(d))
        Convey("各个权重命中是否正确？", func() {
            //各个权重命中是否正确
            So(counter.Get(hosts1[0].Name).Count, ShouldEqual, int64(100*hosts1[0].Weight))
            So(counter.Get(hosts1[1].Name).Count, ShouldEqual, int64(100*hosts1[1].Weight))
            So(counter.Get(hosts1[2].Name).Count, ShouldEqual, int64(100*hosts1[2].Weight))
        })
    })
}

func getTotalWeight(hosts []*HostInstance) int32 {
    c := int32(0)
    for _, v := range hosts {
        c = c + v.Weight
    }
    return c
}
