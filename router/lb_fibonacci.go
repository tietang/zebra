package router

import (
	"github.com/rcrowley/go-metrics"
	"sync"
	"time"
)

type FibonacciWeightRobinRound struct {
	WeightRobinRound
	Base             int64
	weightRobinRound WeightRobinRound
}

func NewFibonacciWeightRobinRound() *FibonacciWeightRobinRound {
	fwrr := &FibonacciWeightRobinRound{}
	fwrr.weightRobinRound = WeightRobinRound{Lock: new(sync.Mutex)}
	fwrr.Lock = new(sync.Mutex)
	fwrr.Base = 100
	return fwrr
}

func (r *FibonacciWeightRobinRound) Next(key string, hosts []*HostInstance) *HostInstance {
	if len(hosts) == 0 {
		return nil
	}
	r.Lock.Lock()
	defer r.Lock.Unlock()

	//保证所有的实例都已经被统计
	ct := 0
	cpHosts := make([]*HostInstance, 0)
	for i := 0; i < len(hosts); i++ {
		hosts[i].Init()
		//  如果有响应时间信息用Fibonacci算法；否则加权轮询算法
		_, r1, _ := r.getResponseTime(hosts[i].InstanceId) // hostsResTime[hosts[i].Name]
		//有实例响应时间，算出权重；否则默认最大权重
		if r1 > 100 {
			ct++
		}
		clone := *hosts[i]
		cpHosts = append(cpHosts, &clone)
	}

	//如果所有的实例被统计，则启用Fibonacci
	if ct == len(hosts) {
		return r.fibonacciNext(key, cpHosts)
	} else {
		return r.weightRobinRound.Next(key, hosts)
	}
}

func (r *FibonacciWeightRobinRound) fibonacciNext(key string, hosts []*HostInstance) *HostInstance {
	selected := hosts[0]
	totalWeight := int32(0)
	for i := 0; i < len(hosts); i++ {
		hosts[i].Init()
		//有实例响应时间，算出权重；否则默认最大权重
		//  如果有响应时间信息用Fibonacci算法；否则加权轮询算法
		res, _, _ := r.getResponseTime(hosts[i].InstanceId) // hostsResTime[hosts[i].Name]
		hosts[i].EffectWeight = int32(GetCeiling(res, r.Base))
		hosts[i].CurrentWeight = hosts[i].EffectWeight
		totalWeight = totalWeight + hosts[i].EffectWeight
		if hosts[i].CurrentWeight > selected.CurrentWeight {
			selected = hosts[i]
		}
	}
	selected.CurrentWeight = selected.CurrentWeight - totalWeight
	return selected
}

func (r *FibonacciWeightRobinRound) getResponseTime(host string) (int64, int64, bool) {
	timer := metrics.GetOrRegisterTimer(host, InstanceRegistry).Snapshot()
	if timer == nil {
		return 1, 0, false
	}
	r1 := int64(timer.Rate1())
	mean := int64(timer.Mean() / float64(time.Millisecond))
	if mean <= 0 {
		return 1, 0, false
	}
	return mean, r1, true
}
