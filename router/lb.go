package router

import (
    "sync"
    "sync/atomic"
    "time"
)

var DefaultBalancerFactory *BalancerFactory

func init() {
    DefaultBalancerFactory = &BalancerFactory{
        defaultFunc: func() Balancer {
            return &WeightRobinRound{Lock: new(sync.Mutex)}
        },
        funcs: make(map[string]func() Balancer),
    }
    DefaultBalancerFactory.Add("WeightRobinRound", func() Balancer {
        return &WeightRobinRound{Lock: new(sync.Mutex)}
    })
    DefaultBalancerFactory.Add("random", func() Balancer {
        return &RandomBalancer{}
    })
    DefaultBalancerFactory.Add("round", func() Balancer {
        return &RoundBalancer{}
    })
    DefaultBalancerFactory.Add("hash", func() Balancer {
        return &HashBalancer{}
    })
    DefaultBalancerFactory.Add("fibonacci", func() Balancer {
        fwrr := NewFibonacciWeightRobinRound()
        //&FibonacciWeightRobinRound{}
        //fwrr.Lock = new(sync.Mutex)
        //fwrr.Base = 100
        return fwrr
    })
}

type Balancer interface {
    Next(key string, hosts []*HostInstance) *HostInstance
}

type InstanceStats struct {
}

type HostInstance struct {
    //负载均衡参数
    isInit        bool
    Weight        int32
    CurrentWeight int32
    EffectWeight  int32
    IsAvailable   bool
    //服务调用
    InstanceId           string
    Name                 string
    AppName              string
    AppGroupName         string
    DataCenter           string
    Tags                 []string
    Params               map[string]string
    ServiceSource        string
    ExternalInstance     interface{}
    Scheme               string
    Address              string
    Port                 string
    HealthCheckUrl       string
    Status               string
    LastUpdatedTimestamp int
    OverriddenStatus     string
    //
    lastSleepOpenTime          *time.Time //最后Sleep打开时间
    lastSleepExpectedCloseTime *time.Time //最后Sleep预计关闭时间
    lastSleepWindowSeqCount    int        //最后连续Sleep统计
    isFailedSleepOpen          bool
    isNotFound                 bool
    //
    lock *sync.Mutex
}

func (h *HostInstance) init() {
    if !h.isInit {
        h.lock = new(sync.Mutex)
        h.EffectWeight = h.Weight
        h.CurrentWeight = h.Weight
        h.isInit = true
    }
}

func (h *HostInstance) ToKey() string {
    return h.InstanceId
}

func (h *HostInstance) SetCurrentWeight(newCurrentWeight int32) {
    //h.lock.Lock()
    //defer h.lock.Unlock()
    //h.CurrentWeight = CurrentWeight
    //atomic.StoreInt32(&h.CurrentWeight, newCurrentWeight)
    atomic.SwapInt32(&h.CurrentWeight, newCurrentWeight)
}

func (h *HostInstance) ResetFailedSleepFlag() {
    h.lock.Lock()
    defer h.lock.Unlock()
    h.lastSleepOpenTime = nil
    h.lastSleepExpectedCloseTime = nil
    h.lastSleepWindowSeqCount = 0
    h.isFailedSleepOpen = false
}

func (h *HostInstance) SetFailedSleepFlagForSeqMode(sleepDuration time.Duration) {
    h.lock.Lock()
    defer h.lock.Unlock()
    lastSleepOpenTime := time.Now()
    //如果多次连续Sleep的第一次，则设置第一次Sleep时间
    if h.lastSleepOpenTime == nil {
        h.lastSleepOpenTime = &lastSleepOpenTime
    } else {
        lastSleepOpenTime = *h.lastSleepOpenTime
    }
    closeTime := lastSleepOpenTime.Add(sleepDuration)
    h.lastSleepExpectedCloseTime = &closeTime
    h.lastSleepWindowSeqCount++
    h.isFailedSleepOpen = true
}

func (h *HostInstance) SetFailedSleepFlagForFixedMode(sleepDuration time.Duration) {
    h.lock.Lock()
    defer h.lock.Unlock()
    lastSleepOpenTime := time.Now()
    h.lastSleepOpenTime = &lastSleepOpenTime
    closeTime := lastSleepOpenTime.Add(sleepDuration)
    h.lastSleepExpectedCloseTime = &closeTime

    h.lastSleepWindowSeqCount++
    h.isFailedSleepOpen = true
}

//复制原生属性
func (h *HostInstance) CopyFrom(source *HostInstance) {
    h.InstanceId = source.InstanceId
    h.Name = source.Name
    h.AppName = source.AppName
    h.AppGroupName = source.AppGroupName
    h.DataCenter = source.DataCenter
    //h.isInit = source.isInit
    if h.EffectWeight == h.Weight {
        h.EffectWeight = source.Weight
    }
    h.Weight = source.Weight
    h.Tags = source.Tags
    h.Params = source.Params
    h.ServiceSource = source.ServiceSource
    h.ExternalInstance = source.ExternalInstance
    h.Scheme = source.Scheme
    h.Address = source.Address
    h.Port = source.Port
    h.HealthCheckUrl = source.HealthCheckUrl
    h.Status = source.Status
    h.LastUpdatedTimestamp = source.LastUpdatedTimestamp
    h.OverriddenStatus = source.OverriddenStatus
}

//片排序
type ByCurrentWeight []*HostInstance

func (p ByCurrentWeight) Len() int {
    return len(p)
}

func (p ByCurrentWeight) Less(i, j int) bool {
    return p[i].CurrentWeight >= p[j].CurrentWeight
}

func (p ByCurrentWeight) Swap(i, j int) {
    p[i], p[j] = p[j], p[i]
}

type BalancerFactory struct {
    defaultFunc func() Balancer
    funcs       map[string]func() Balancer
}

func (bf *BalancerFactory) Add(lbName string, fun func() Balancer) {
    bf.funcs[lbName] = fun
}
func (bf *BalancerFactory) New(lbName string) Balancer {
    if fun, ok := bf.funcs[lbName]; ok {
        return fun()
    }
    return bf.defaultFunc()

}
