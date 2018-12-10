package router

import (
    log "github.com/sirupsen/logrus"
    "github.com/tietang/go-eureka-client/eureka"
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/discovery"
    "github.com/tietang/zebra/health"
    "strconv"
    "strings"
    "time"
)

const (
    DEFAULT_EUREKA_SERVER_URLS = "http://127.0.0.1:8761/eureka"
    //key
    KEY_EUREKA_ENABLED            = "eureka.server.enabled"
    KEY_EUREKA_SERVER_URLS        = "eureka.server.urls"
    KEY_EUREKA_DISCOVERY_INTERVAL = "eureka.discovery.interval"
    //default value
    DEFAULT_EUREKA_SERVER_ENABLED     = "false"
    DEFAULT_EUREKA_DISCOVERY_INTERVAL = time.Duration(10000) * time.Millisecond
)

type EurekaRouteSource struct {
    KeyValueRouteSource
    urls      []string
    discovery *discovery.EurekaDiscovery
}

func NewEurekaRouteSource(conf kvs.ConfigSource) *EurekaRouteSource {
    e := &EurekaRouteSource{}
    e.conf = conf
    urls := e.conf.GetDefault(KEY_EUREKA_SERVER_URLS, DEFAULT_EUREKA_SERVER_URLS)

    eurekaUrls := strings.Split(urls, ",|, | , ")
    e.name = "eureka:" + urls
    e.urls = eurekaUrls
    //eurekaUrl []string, router *Router, conf kvs.conf

    return e
}

func (p *EurekaRouteSource) Init() {
    if p.isInit {
        return
    }

    discovery := discovery.NewEurekaDiscovery(p.urls)
    discovery.ScheduleAtFixedRate(p.conf.GetDurationDefault(KEY_EUREKA_DISCOVERY_INTERVAL, DISCOVERY_INTERVAL_DEFAULT))
    p.discovery = discovery
    p.updateDiscoveryRouter(discovery.GetApps())
    p.discovery.AddCallback(p.updateDiscoveryRouter)
    p.isInit = true
}

func (h *EurekaRouteSource) updateDiscoveryRouter(apps *eureka.Applications) {
    if apps == nil || apps.Applications == nil {
        log.Debug("discovery apps is empty")
        return
    }
    log.Debug("update router by eureka.")
    for _, a := range apps.Applications {
        appName := strings.ToLower(a.Name)
        route := &Route{
            Source:      "/" + appName + "/**",
            ServiceId:   a.Name,
            StripPrefix: true,
        }

        h.AddInTime(route)

        //log.WithField("route", route).Debug("add router by eureka: ")
        size := len(a.Instances)
        hosts := make([]*HostInstance, size)
        for i, instance := range a.Instances {
            //h.hostChangedCallback(appName, newHostInstanceByEureka(appName, &instance))
            hosts[i] = newHostInstanceByEureka(appName, &instance)
        }
        h.appHostsChangedCallback(appName, hosts)

    }
}

func newHostInstanceByEureka(appName string, ins *eureka.InstanceInfo) *HostInstance {

    port := ins.Port.Port
    scheme := "http"
    if ins.SecurePort.Enabled {
        port = ins.SecurePort.Port
        scheme = "https"
    }
    weight := int32(1)
    //if weightStr, ok := ins.Metadata.Map["weight"]; ok {
    //    w, err := strconv.Atoi(weightStr)
    //    if err == nil {
    //        weight = w
    //    }
    //}

    h := &HostInstance{
        AppGroupName:         ins.AppGroupName,
        AppName:              appName,
        DataCenter:           ins.DataCenterInfo.Name,
        Address:              ins.IpAddr,
        Port:                 strconv.Itoa(port),
        InstanceId:           ServiceId(appName, ins.IpAddr, strconv.Itoa(port)),
        Name:                 strings.Join([]string{ins.HostName, strconv.Itoa(port)}, ":"),
        LastUpdatedTimestamp: ins.LastUpdatedTimestamp,
        Status:               ins.Status,
        OverriddenStatus:     ins.Overriddenstatus,
        HealthCheckUrl:       ins.HealthCheckUrl,
        Scheme:               scheme,
        ServiceSource:        "eureka",
        ExternalInstance:     ins,
        Params:               ins.Metadata.Map,
        Weight:               weight,
    }
    h.init()
    return h
}

func (h *EurekaRouteSource) CheckHealth(rootHealth *health.RootHealth) {
    if h.discovery != nil {
        //http://172.16.1.248:8500/v1/status/leader
        ok, desc := h.discovery.Health()

        if ok {
            rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_UP, Desc: desc}

        } else {
            rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_DOWN, Desc: desc}
            rootHealth.Status = health.STATUS_DOWN
        }
    }
}
