package discovery

import (
    "github.com/hashicorp/consul/api"
    log "github.com/sirupsen/logrus"
    "strings"
    "time"
)

type ConsulDiscovery struct {
    Config          *api.Config
    catalogServices map[string][]*api.CatalogService
    services        map[string][]string
    callbacks       []func(map[string][]string, map[string][]*api.CatalogService)
    client          *api.Client
}

func NewConsulDiscoveryByConfig(config *api.Config) *ConsulDiscovery {
    if config == nil {
        config = api.DefaultConfig()
    }
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }

    cd := &ConsulDiscovery{Config: config, callbacks: make([]func(map[string][]string, map[string][]*api.CatalogService), 0)}
    cd.client = client
    return cd
}

//address: [hostName:port]
func NewConsulDiscovery(address string) *ConsulDiscovery {
    config := api.DefaultConfig()
    config.Address = address
    return NewConsulDiscoveryByConfig(config)
}

func (d *ConsulDiscovery) AddCallback(callback func(map[string][]string, map[string][]*api.CatalogService)) {
    d.callbacks = append(d.callbacks, callback)

}
func (c *ConsulDiscovery) GetServicesInTime() (map[string][]string, map[string][]*api.CatalogService, error) {

    q := &api.QueryOptions{}
    services, _, err := c.client.Catalog().Services(q)
    if err != nil {
        return nil, nil, err
    }

    catalogServices := make(map[string][]*api.CatalogService)

    for serviceName, _ := range services {
        cs, _, err := c.client.Catalog().Service(serviceName, "", q)
        if err != nil {
            log.Error(err)
            continue
        }
        catalogServices[serviceName] = cs

    }

    return services, catalogServices, nil

}

func (c *ConsulDiscovery) GetServices() (map[string][]string, map[string][]*api.CatalogService) {
    if c.catalogServices == nil {
        services, catalogServices, err := c.GetServicesInTime()
        if err == nil {
            return services, catalogServices
        }
    }
    return c.services, c.catalogServices
}

func (c *ConsulDiscovery) GetService(name string) []*api.CatalogService {
    if c.catalogServices == nil {
        log.Info("catalogServices is nil")
        return nil
    }
    for name, service := range c.catalogServices {
        if strings.ToLower(name) == strings.ToLower(name) {
            return service
        }
    }
    return nil
}

func (d *ConsulDiscovery) ScheduleAtFixedRate(second time.Duration) {
    d.run()
    go d.runTask(second)
}

func (d *ConsulDiscovery) runTask(second time.Duration) {
    timer := time.NewTicker(second)
    for {
        select {
        case <-timer.C:
            go d.run()
        }
    }
}
func (d *ConsulDiscovery) run() {
    services, catalogServices, err := d.GetServicesInTime()
    if err == nil || services != nil || catalogServices != nil {

        //for key, _ := range d.services {
        //    hasExists := false
        //    for name, _ := range services {
        //        if name == key {
        //            hasExists = true
        //        }
        //    }
        //    if !hasExists {
        //        delete(services, key)
        //    }
        //}
        d.services = services

        //for key, _ := range d.catalogServices {
        //    hasExists := false
        //    for name, _ := range catalogServices {
        //        if name == key {
        //            hasExists = true
        //        }
        //    }
        //    if !hasExists {
        //        delete(catalogServices, key)
        //    }
        //}

        d.catalogServices = catalogServices
        d.execCallbacks(d.services, d.catalogServices)
    } else {
        log.Error(err)
    }
}

func (d *ConsulDiscovery) execCallbacks(services map[string][]string, catalogServices map[string][]*api.CatalogService) {
    if len(d.callbacks) > 0 {
        for _, c := range d.callbacks {
            go c(services, catalogServices)
        }
    }
}

func (c *ConsulDiscovery) Health() (bool, string) {
    leader, err := c.client.Status().Leader()
    if err != nil || leader == "" {
        return false, err.Error()
    }

    ok, desc := true, "ok"

    return ok, desc

}
