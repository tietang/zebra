package router

import (
	"github.com/hashicorp/consul/api"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/consul"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/discovery"
	"github.com/tietang/zebra/health"
	"strconv"
	"strings"
	"time"
)

const (
	KEY_CONSUL_ENABLED            = "consul.enabled"
	KEY_CONSUL_DISCOVERY_ADDRESS  = "consul.discovery.address"
	KEY_CONSUL_DISCOVERY_ENABLED  = "consul.discovery.enabled"
	KEY_CONSUL_DISCOVERY_INTERVAL = "consul.discovery.interval"

	KEY_CONSUL_ROUTES_ROOT    = "consul.routes.root"
	KEY_CONSUL_ROUTES_ADDRESS = "consul.address"
	KEY_CONSUL_ROUTES_ENABLED = "consul.routes.enabled"

	DEFAULT_CONSUL_ROUTES_ROOT = "zebra/routes"
	DEFAULT_CONSUL_ADDRESS     = "127.0.0.1:8500"
)

type ConsulRouteSource struct {
	KeyValueRouteSource
	address            string
	discovery          *discovery.ConsulDiscovery
	consulConfigSource *consul.ConsulKeyValueConfigSource
}

func NewConsulRouteSource(conf kvs.ConfigSource) *ConsulRouteSource {
	e := &ConsulRouteSource{}
	e.conf = conf
	//eurekaUrl []string, router *Router, conf kvs.conf
	e.address = conf.GetDefault(KEY_CONSUL_DISCOVERY_ADDRESS, DEFAULT_CONSUL_ADDRESS)
	e.name = "consul:" + e.address
	return e
}

func (p *ConsulRouteSource) Init() {
	if p.isInit {
		return
	}

	if p.conf.GetBoolDefault(KEY_CONSUL_DISCOVERY_ENABLED, false) {
		log.Info("consul discovery enabled.")
		discovery := discovery.NewConsulDiscovery(p.address)
		discovery.ScheduleAtFixedRate(p.conf.GetDurationDefault(KEY_CONSUL_DISCOVERY_INTERVAL, DISCOVERY_INTERVAL_DEFAULT))
		p.discovery = discovery
		p.updateDiscoveryRouter(p.discovery.GetServices())
		p.discovery.AddCallback(p.updateDiscoveryRouter)

	}
	//

	if p.conf.GetBoolDefault(KEY_CONSUL_ROUTES_ENABLED, false) {
		log.Info("consul routes enabled.")
		root := p.conf.GetDefault(KEY_CONSUL_ROUTES_ROOT, DEFAULT_CONSUL_ROUTES_ROOT)
		address := p.conf.GetDefault(KEY_CONSUL_ROUTES_ADDRESS, DEFAULT_CONSUL_ADDRESS)

		source := consul.NewConsulKeyValueConfigSource(address, root)

		p.consulConfigSource = source
		p.loadRouteByConfigSource(p.consulConfigSource)

	}

	p.isInit = true
}
func getRoutePrefix(tags []string) string {
	for _, tag := range tags {
		if strings.Contains(tag, KEY_LABEL_ROUTE_PREFIX) {
			kv := strings.Split(tag, "=")
			if len(kv) >= 2 {
				return kv[1]
			}
		}
	}
	return ""
}
func (h *ConsulRouteSource) updateDiscoveryRouter(services map[string][]string, catalogServices map[string][]*api.CatalogService) {
	//log.Info(services)
	//log.Info(catalogServices)
	if services == nil {
		log.Debug("consul services is empty")
		return
	}
	for name, tags := range services {
		appName := strings.ToLower(name)
		routePrefix := getRoutePrefix(tags)

		urlPattern := toPath(routePrefix)

		if routePrefix == "" {
			urlPattern = "/" + appName + "/**"
		}

		route := &Route{
			Source:      urlPattern,
			ServiceId:   name,
			StripPrefix: true,
		}

		h.Add(route)
		log.WithField("route", route).Debug("add router by consul: ")
		for appName, instances := range catalogServices {
			//for _, instance := range instances {
			//    h.hostChangedCallback(appName, newHostInstanceByConsul(appName, instance))
			//
			//}

			size := len(instances)
			hosts := make([]*HostInstance, size)
			for i, instance := range instances {
				//h.hostChangedCallback(appName, newHostInstanceByEureka(appName, &instance))
				hosts[i] = newHostInstanceByConsul(appName, instance)
			}
			h.appHostsChangedCallback(appName, hosts)

		}

	}
}

func newHostInstanceByConsul(appName string, ins *api.CatalogService) *HostInstance {

	port := ins.ServicePort
	scheme := "http"

	version := ins.ServiceMeta["version"]
	h := &HostInstance{
		AppGroupName:         ins.ServiceName,
		AppName:              ins.ServiceName,
		DataCenter:           ins.Datacenter,
		Address:              ins.ServiceAddress,
		Port:                 strconv.Itoa(port),
		InstanceId:           ServiceId(ins.ServiceName, ins.ServiceAddress, strconv.Itoa(ins.ServicePort)),
		Name:                 strings.Join([]string{ins.ServiceAddress, strconv.Itoa(port)}, ":"),
		LastUpdatedTimestamp: int(time.Now().UnixNano() / int64(time.Millisecond)),
		Status:               STATUS_UP,
		//OverriddenStatus:     "",
		//HealthCheckUrl:       ins.HealthCheckUrl,
		Scheme:           scheme,
		ServiceSource:    "consul",
		ExternalInstance: ins,
		Tags:             ins.ServiceTags,
		//Labels:           ins.Metadata.Map,
		Weight:  1,
		Version: version,
	}
	h.init()
	return h
}
func (h *ConsulRouteSource) CheckHealth(rootHealth *health.RootHealth) {
	if h.discovery != nil {
		ok, desc := h.discovery.Health()
		if ok {
			rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_UP, Desc: desc}

		} else {
			rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_DOWN, Desc: desc}
			rootHealth.Status = health.STATUS_DOWN
		}
	}
}

//
//func hasTag(tags []string, tag string) (bool, string) {
//	for _, name := range tags {
//		if tag == name {
//			return true
//		}
//	}
//	return false
//}
