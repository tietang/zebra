package router

import (
	"sync"
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
