package discovery

import (
    "encoding/json"
    "github.com/hashicorp/consul/api"
    log "github.com/sirupsen/logrus"
    "path"
    "strings"
    "sync"
    "time"
)

// 通过consul kv store来服务发现，通过外部配置等等
// 配置规则：
// 1. 指定一个root path作为key prefix
// 2. key prefix之后的一级key作为service name
// 3. key prefix之后的第二级key作为服务实例，格式为：hostName|IP:port
// 举例：
// key prefix: example/services/
// service name: customer
// 服务实例：192.168.1.2:8080,192.168.1.3:8080,192.168.1.4:8080
// 那么：
//http://172.16.1.2:8500/v1/kv/example/services?keys
//[
//"example/services/customer/192.168.1.2:8080",
//"example/services/customer/192.168.1.3:8080",
//"example/services/customer/192.168.1.4:8080"
//]
//value：
// {
//"InstanceId": "customer:192.168.1.4:8080",
//"Name": "192.168.1.4:8080",
//"AppName": "customer",
//"AppGroupName": "customer",
//"Tags": null,
//"Labels": null,
//"InstanceType": "consul_kv",
//"ExternalInstance": null,
//"Scheme": "http",
//"Address": "192.168.1.4",
//"Port": "8080",
//"HealthCheckUrl": "/health",
//"Status": "UP",
//"LastUpdatedTimestamp": 0,
//"OverriddenStatus": ""
//}
type ConsulKeyValueDiscovery struct {
    Config    *api.Config
    services  map[string]*Service
    callbacks []func(map[string]*Service)
    client    *api.Client
    kv        *api.KV
    root      string
    lock      *sync.Mutex
}

func NewConsulKeyValueDiscoveryByConfig(config *api.Config, root string) *ConsulKeyValueDiscovery {
    if config == nil {
        config = api.DefaultConfig()
    }
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }
    cd := &ConsulKeyValueDiscovery{Config: config, callbacks: make([]func(map[string]*Service), 0)}
    cd.client = client
    cd.lock = new(sync.Mutex)
    cd.root = root
    cd.kv = client.KV()
    return cd
}

//address: [hostName:port]
func NewConsulKeyValueDiscovery(address string, root string) *ConsulKeyValueDiscovery {
    config := api.DefaultConfig()
    config.Address = address
    return NewConsulKeyValueDiscoveryByConfig(config, root)
}

func (d *ConsulKeyValueDiscovery) AddCallback(callback func(map[string]*Service)) {
    d.callbacks = append(d.callbacks, callback)

}
func (c *ConsulKeyValueDiscovery) GetServicesInTime() (map[string]*Service, error) {

    q := &api.QueryOptions{}
    prefix := c.root
    keys, _, err := c.kv.Keys(prefix, "", q)
    if err != nil {
        return nil, err
    }
    services := make(map[string]*Service, 0)

    for _, k := range keys {
        kv, _, err := c.kv.Get(k, q)
        if err != nil {
            continue
        }
        name := path.Base(k)
        ins := &Instance{}
        err = json.Unmarshal(kv.Value, ins)
        if err != nil {
            continue
        }
        c.add(services, name, ins)

    }

    return services, nil

}

func (c *ConsulKeyValueDiscovery) add(services map[string]*Service, name string, ins *Instance) {
    c.lock.Lock()
    defer c.lock.Unlock()
    if service, ok := services[name]; ok {
        service.Instances = append(service.Instances, ins)
    } else {
        service := &Service{
            Name:      name,
            Instances: make([]*Instance, 0),
        }
        services[name] = service
    }

}

func (c *ConsulKeyValueDiscovery) GetServices() (map[string]*Service) {
    if c.services == nil {
        services, err := c.GetServicesInTime()
        if err == nil {
            return services
        }
    }
    return c.services
}

func (c *ConsulKeyValueDiscovery) GetService(name string) *Service {
    if c.services == nil {
        log.Info("catalogServices is nil")
        return nil
    }
    for name, service := range c.services {
        if strings.ToLower(name) == strings.ToLower(name) {
            return service
        }
    }
    return nil
}

func (d *ConsulKeyValueDiscovery) Watching(second time.Duration) {
    d.run()
    go d.runTask(second)
}

func (d *ConsulKeyValueDiscovery) runTask(second time.Duration) {
    timer := time.NewTicker(second)
    for {
        select {
        case <-timer.C:
            go d.run()
        }
    }
}
func (d *ConsulKeyValueDiscovery) run() {
    services, err := d.GetServicesInTime()
    if err == nil || services != nil {
        d.services = services
        d.execCallbacks(d.services)
    } else {
        log.Error(err)
    }
}

func (d *ConsulKeyValueDiscovery) execCallbacks(services map[string]*Service) {
    if len(d.callbacks) > 0 {
        for _, c := range d.callbacks {
            go c(services)
        }
    }
}

func (c *ConsulKeyValueDiscovery) Health() (bool, string) {
    leader, err := c.client.Status().Leader()
    if err != nil || leader == "" {
        return false, err.Error()
    }

    ok, desc := true, "ok"

    return ok, desc

}
