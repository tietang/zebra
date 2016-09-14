package zuul

import (
    "strconv"
    "strings"

    "github.com/tietang/go-eureka-client/eureka"
)

type DiscoveryRobin struct {
    Balancer  Balancer
    discovery *Discovery
    hosts     map[string][]HostInstance
}

func NewDiscoveryRobin(balancer Balancer, discovery *Discovery) *DiscoveryRobin {
    return &DiscoveryRobin{Balancer: balancer, discovery: discovery, hosts: make(map[string][]HostInstance)}
}

func (d *DiscoveryRobin) Next(appName string) (*eureka.Application, *HostInstance, *eureka.InstanceInfo) {

    var app eureka.Application
    apps := d.discovery.apps

    if apps == nil || apps.Applications == nil {
        return nil, nil, nil
    }
    for _, a := range apps.Applications {
        if strings.ToUpper(appName) == strings.ToUpper(a.Name) {
            app = a
            if len(a.Instances) == 0 {
                return &app, nil, nil
            }
            for _, ins := range a.Instances {
                if ins.Status == eureka.UP {
                    name := d.ins2name(&ins)
                    host := HostInstance{Id: len(d.hosts[appName]), Name: name, Weight: 1}
                    if !d.hasExists(appName, host) {
                        d.hosts[appName] = append(d.hosts[appName], host)
                    }
                }
            }
        }
    }

    host := d.Balancer.Next(d.hosts[appName])
    ins := d.findIns(app.Instances, host)
    return &app, host, ins

}
func (d *DiscoveryRobin) findIns(instances []eureka.InstanceInfo, host *HostInstance) *eureka.InstanceInfo {
    for _, ins := range instances {
        name := d.ins2name(&ins)
        if name == host.Name {
            return &ins
        }
    }
    return nil
}

func (d *DiscoveryRobin) ins2name(ins *eureka.InstanceInfo) string {
    var port int
    if ins.Port.Enabled {
        port = ins.Port.Port
    } else {
        port = ins.SecurePort.Port
    }
    name := ins.IpAddr + ":" + strconv.Itoa(port)

    return name
}

func (d *DiscoveryRobin) hasExists(appName string, hosts HostInstance) bool {
    for _, v := range d.hosts[appName] {
        if v.Name == hosts.Name {
            return true
        }
    }
    return false
}
