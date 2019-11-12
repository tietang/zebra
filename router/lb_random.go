package router

import (
    "math/rand"
)

//type HostInstances []HostInstance

type RandomBalancer struct {
}

func (r *RandomBalancer) Next(key string, hosts []*HostInstance) *HostInstance {
    size := len(hosts)
    seed := rand.Uint64()
    index := int(seed) % size
    if index < 0 {
        index = -index
    }
    selected := hosts[index]
    return selected
}
