package router

import (
    "github.com/lafikl/consistent"
    log "github.com/sirupsen/logrus"
)

type HashBalancer struct {
}

func (r *HashBalancer) Next(key string, hosts []*HostInstance) *HostInstance {
    c := consistent.New()
    for _, ins := range hosts {
        c.Add(ins.Name)
    }
    hostName, err := c.Get(key)
    if err != nil {
        log.Fatal(err)
    }
    for _, ins := range hosts {
        if hostName == ins.Name {
            return ins
        }
    }

    return nil
}
