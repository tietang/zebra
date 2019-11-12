package router

import (
    "github.com/go-ini/ini"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/health"
    "reflect"
    "strconv"
    "strings"
    "sync/atomic"
)

var globalRouterSources []RouterSource

func init() {
    globalRouterSources = make([]RouterSource, 0)
}

type HostInstanceSource interface {
    Name() string
    //开始初始化和监听
    Start()
    //设置变更通知回调函数
    SetAppHostsChangedCallback(callback func(appName string, h []*HostInstance))
}

type RouteSource interface {
    Name() string
    //设置变更通知回调函数
    SetRouterChangedCallback(callback func(route *Route))
    //初始化：连接，初始信息构建
    Init()
    //开始初始化和监听
    Start()
    //检查提供该服务的健康状态
    CheckHealth(rootHealth *health.RootHealth)
}

type RouterSource interface {
    Build()
    Name() string
    SetRouterChangedCallback(callback func(route *Route))
    Init()
    Start()
    CheckHealth(rootHealth *health.RootHealth)
    SetAppHostsChangedCallback(callback func(appName string, h []*HostInstance))
}

//自定义扩展route和hostinstance扩展注册：
//实现RouterSource的所有方法
func RegisterSource(rs RouterSource) {
    globalRouterSources = append(globalRouterSources, rs)
}

type KeyValueRouteSource struct {
    name                    string
    isEnabled               bool
    isInit                  bool
    isStarted               bool
    routerStripPrefix       bool
    routes                  []*Route
    routeChangedCallback    func(route *Route)
    appHostsChangedCallback func(appName string, h []*HostInstance)
    conf                    kvs.ConfigSource
    StripPrefix             bool
}

var keyValueRouteSourcePosition int32 = 0

func NewKeyValueRouteSource() *KeyValueRouteSource {
    atomic.AddInt32(&keyValueRouteSourcePosition, 1)
    k := &KeyValueRouteSource{
        name:   "KeyValueRouteSource-" + strconv.Itoa(int(keyValueRouteSourcePosition)),
        routes: make([]*Route, 0),
        conf:   kvs.NewEmptyNoSystemEnvCompositeConfigSource(),
    }

    return k
}

func (p *KeyValueRouteSource) Init() {
    if p.isInit {
        return
    }
    p.isInit = true
}
func (p *KeyValueRouteSource) Add(route *Route) {
    p.add(route, false)
}
func (p *KeyValueRouteSource) AddInTime(route *Route) {
    p.add(route, true)
}
func (p *KeyValueRouteSource) add(route *Route, isInTimeSync bool) {
    if route == nil {
        return
    }
    route.Init()
    hasExists := false
    for i, r := range p.routes {
        if r.Id == route.Id {
            hasExists = true
            p.routes[i] = route
            log.WithField("route", route).Debug("update route: ")
        }
        //对于通过erueka，k8s,consul,zk discovery等自动服务发现的规则，如果设置了IsForceUpdate=true，可以动态更新。
        if route.IsForceUpdate && route.ServiceId == r.ServiceId {
            hasExists = true
            p.routes[i] = route
            log.WithField("route", route).Debug("force update route: ")
        }
    }
    if !hasExists {
        p.routes = append(p.routes, route)
        log.WithField("route", route).Debug("add route: ")
    }
    if isInTimeSync && p.routeChangedCallback != nil {
        p.routeChangedCallback(route)
    }
    if p.routeChangedCallback == nil {
        log.Warn("RouteSource routeChangedCallback is not given.")
    }
    configHystrix(route.ServiceId, p.conf)
}

func (p *KeyValueRouteSource) Name() string {
    return p.name
}

func (p *KeyValueRouteSource) Start() {
    p.Init()
    if p.isStarted {
        return
    }
    for _, route := range p.routes {
        if p.routeChangedCallback != nil {
            p.routeChangedCallback(route)
        } else {
            log.Warn("RouteSource routeChangedCallback is not given.")
        }
    }
    p.isStarted = true
}

func (p *KeyValueRouteSource) SetRouterChangedCallback(callback func(route *Route)) {
    p.routeChangedCallback = callback
}

func (p *KeyValueRouteSource) SetAppHostsChangedCallback(callback func(appName string, h []*HostInstance)) {
    p.appHostsChangedCallback = callback
}

func (p *KeyValueRouteSource) addConfigSource(source kvs.ConfigSource) {
    confValue := reflect.ValueOf(p.conf)
    typ := confValue.Elem().Type()
    if typ.String() == "kvs.CompositeConfigSource" {
        conf := p.conf.(*kvs.CompositeConfigSource)
        conf.Add(source)
    } else {
        conf := kvs.NewEmptyNoSystemEnvCompositeConfigSource()
        //NewCompositeConfigSource("IniFileCompositeConfigSource", p.conf, source)
        conf.Add(p.conf, source)
        p.conf = conf
    }
}

func (g *KeyValueRouteSource) loadRouteByConfigSource(source kvs.ConfigSource) {
    keys := source.Keys()
    log.Info("add router by config: ", keys)
    log.Info("by source: ", source.Name())
    for _, k := range keys {

        //if strings.HasPrefix(k, KEY_ROUTES_NODE_PREFIX) {

        v, err := source.Get(k)
        if err != nil {
            continue
        }
        //log.Info("append blank: ", v)
        ///app1,app1/v1/user,app1,/v1/user/app1/info,app1,/info
        //source path, app name,  target path
        //
        props, err := ini.Load([]byte(v))
        if err != nil {
            log.Warn(err, ": ", k)
            continue
        }
        sections := props.Sections()

        sp := readIniSections(g.routerStripPrefix, sections, func(route *Route) {
            g.Add(route)
        })
        g.addConfigSource(sp)
    }

    g.Start()

}

func ServiceId(name, ip, portStr string) string {
    return strings.Join([]string{name, ip, portStr}, ":")
}
