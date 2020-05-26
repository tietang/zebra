package router

import (
	"sync"
)

//type HostInstances []HostInstance

type WeightRobinRound struct {
	Lock *sync.Mutex
}

func (r *WeightRobinRound) Next(key string, hosts []*HostInstance) *HostInstance {
	if len(hosts) == 0 {
		return nil
	}
	r.Lock.Lock()
	defer r.Lock.Unlock()
	selected := hosts[0]
	totalWeight := int32(0)
	for i := 0; i < len(hosts); i++ {
		hosts[i].Init()
		totalWeight = totalWeight + hosts[i].EffectWeight
		hosts[i].CurrentWeight = hosts[i].CurrentWeight + hosts[i].EffectWeight
		if hosts[i].CurrentWeight > selected.CurrentWeight {
			selected = hosts[i]
		}
	}
	selected.CurrentWeight = selected.CurrentWeight - totalWeight
	return selected
}
