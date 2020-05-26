package router

//
//import (
//	log "github.com/sirupsen/logrus"
//	"github.com/tietang/props/kvs"
//	"github.com/tietang/zebra/discovery"
//	"github.com/tietang/zebra/health"
//	"strings"
//)
//
////go get k8s.io/client-go/...
//const (
//	DISCOVERY_TYPE = "k8s"
//
//	KEY_K8S_ENABLED            = "k8s.enabled"
//	KEY_K8S_URLS               = "k8s.urls"
//	KEY_K8S_DISCOVERY_INTERVAL = "k8s.discovery.interval"
//
//	KEY_K8S_LABEL_HEALTH_CHECK_URL = "healthCheckUrl"
//	DEFAULT_K8S_URLS               = "http://127.0.0.1:8080"
//)
//
//type KubernetesRouteSource struct {
//	KeyValueRouteSource
//	url       string
//	discovery *discovery.KubernetesDiscovery
//}
//
//func NewKubernetesRouteSource(conf kvs.ConfigSource) *KubernetesRouteSource {
//	e := &KubernetesRouteSource{}
//	e.conf = conf
//	e.url = e.conf.GetDefault(KEY_K8S_URLS, DEFAULT_K8S_URLS)
//
//	e.name = "Kubernetes:" + e.url
//	//kubernetesUrl []string, router *Router, conf kvs.conf
//
//	return e
//}
//
//func (p *KubernetesRouteSource) Init() {
//	if p.isInit {
//		return
//	}
//
//	discovery := discovery.NewKubernetesDiscovery(p.url)
//	discovery.Watching(p.conf.GetDurationDefault(KEY_K8S_DISCOVERY_INTERVAL, DISCOVERY_INTERVAL_DEFAULT))
//	p.discovery = discovery
//	p.initDiscoveryRouter(discovery.GetServices())
//	p.discovery.AddCallback(p.initDiscoveryRouter)
//	p.isInit = true
//}
//
//func (h *KubernetesRouteSource) initDiscoveryRouter(services map[string]*discovery.Service) {
//	if services == nil {
//		log.Debug("Kubernetes apps is empty")
//		return
//	}
//	for _, s := range services {
//		appName := strings.ToLower(s.Name)
//		routePrefix := s.Labels[KEY_LABEL_ROUTE_PREFIX]
//
//		urlPattern := toPath(routePrefix)
//
//		if routePrefix == "" {
//			urlPattern = "/" + appName + "/**"
//		}
//		host := &Route{
//			Source:        urlPattern,
//			ServiceId:     s.Name,
//			StripPrefix:   true,
//			ServiceSource: DISCOVERY_TYPE,
//			IsForceUpdate: true,
//		}
//
//		h.Add(host)
//
//		log.WithField("host", host).Debug("add router by Kubernetes: ")
//
//		//for _, instance := range s.Instances {
//		//    h.hostChangedCallback(appName, newHostInstanceByInstance(appName, instance))
//		//}
//
//		size := len(s.Instances)
//		hosts := make([]*HostInstance, size)
//		for i, instance := range s.Instances {
//			//h.hostChangedCallback(appName, newHostInstanceByEureka(appName, &instance))
//			hosts[i] = newHostInstanceByInstance(appName, instance)
//		}
//		h.appHostsChangedCallback(appName, hosts)
//
//	}
//}
//
//func newHostInstanceByInstance(appName string, ins *discovery.Instance) *HostInstance {
//
//	port := ins.Port
//	scheme := "http"
//	if ins.Scheme != "" {
//		scheme = ins.Scheme
//	}
//
//	healthCheckUrl := ins.HealthCheckUrl
//	if healthCheckUrl == "" {
//		healthCheckUrl = ins.Params[KEY_K8S_LABEL_HEALTH_CHECK_URL]
//	}
//
//	h := &HostInstance{
//		AppGroupName:         ins.AppGroupName,
//		AppName:              appName,
//		DataCenter:           ins.AppGroupName,
//		Address:              ins.Address,
//		Port:                 port,
//		InstanceId:           ServiceId(appName, ins.Address, ins.Port),
//		Name:                 strings.Join([]string{ins.Address, port}, ":"),
//		LastUpdatedTimestamp: ins.LastUpdatedTimestamp,
//		Status:               ins.Status,
//		OverriddenStatus:     ins.OverriddenStatus,
//		HealthCheckUrl:       ins.HealthCheckUrl,
//		Scheme:               scheme,
//		ServiceSource:        DISCOVERY_TYPE,
//		ExternalInstance:     ins,
//		Params:               ins.Params,
//		Tags:                 ins.Tags,
//		Weight:               1,
//	}
//	h.init()
//	return h
//}
//
//func (h *KubernetesRouteSource) CheckHealth(rootHealth *health.RootHealth) {
//	if h.discovery != nil {
//		//http://172.16.1.248:8500/v1/status/leader
//		ok, desc := h.discovery.Health()
//
//		if ok {
//			rootHealth.Healths["kubernetes server"] = &health.Health{Status: health.STATUS_UP, Desc: desc}
//
//		} else {
//			rootHealth.Healths["kubernetes server"] = &health.Health{Status: health.STATUS_DOWN, Desc: desc}
//			rootHealth.Status = health.STATUS_DOWN
//		}
//	}
//}
