package router

import (
	"sync/atomic"
)

//type HostInstances []HostInstance

type RoundBalancer struct {
	lastCt uint64
}

func (r *RoundBalancer) Next(key string, hosts []*HostInstance) *HostInstance {
	new := atomic.AddUint64(&r.lastCt, 1)
	size := len(hosts)
	index := int(new) % size
	selected := hosts[index]
	return selected
}
