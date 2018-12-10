package router

import (
    "github.com/deckarep/golang-set"
    _ "github.com/go-sql-driver/mysql"
    "github.com/go-xorm/xorm"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/go-utils"
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/health"
    "time"
)

const (
    SQL_DRIVERNAME            = "sql.driverName"
    SQL_URL                   = "sql.url"
    KEY_SQL_ROUTES_ENABLED    = "sql.routes.enabled"
    KEY_SQL_DISCOVERY_ENABLED = "sql.discovery.enabled"
)

type RouteModel struct {
    IdRoute            uint64 `xorm:"bigint(20) notnull unique 'id' index 'id' pk autoincr"`
    Id                 string
    ServiceId          string
    Source             string
    SourcePrefix       string
    SourceIsFuzzyMatch bool
    Target             string
    TargetIsFuzzyMatch bool
    TargetPrefix       string
    StripPrefix        bool
    Subdomain          string // "admin."
    Path               string // "/api/user/:id"
    Enabled            bool
    Sorted             int
    HostInstances      []*HostInstanceModel
}

func (r *RouteModel) toRoute() *Route {
    route := &Route{
        Id:                 r.Id,
        ServiceId:          r.ServiceId,
        Source:             r.Source,
        SourcePrefix:       r.SourcePrefix,
        SourceIsFuzzyMatch: r.SourceIsFuzzyMatch,
        Target:             r.Target,
        TargetIsFuzzyMatch: r.TargetIsFuzzyMatch,
        TargetPrefix:       r.TargetPrefix,
        StripPrefix:        r.StripPrefix,
        Subdomain:          r.Subdomain,
        Path:               r.Path,
    }
    return route
}

type HostInstanceModel struct {
    Id                   uint64
    IdRoute              uint64
    InstanceId           string
    Name                 string
    AppName              string
    AppGroupName         string
    DataCenter           string
    Weight               int32
    CurrentWeight        int32
    EffectWeight         int32
    InstanceType         string
    Scheme               string
    Address              string
    Port                 string
    HealthCheckUrl       string
    Status               string
    LastUpdatedTimestamp int
    OverriddenStatus     string
    Enabled              bool
}

func (r *HostInstanceModel) toHostInstance() *HostInstance {
    ins := &HostInstance{
        InstanceId:           r.InstanceId,
        Name:                 r.Name,
        AppName:              r.AppName,
        DataCenter:           r.DataCenter,
        AppGroupName:         r.AppGroupName,
        Weight:               r.Weight,
        CurrentWeight:        r.CurrentWeight,
        EffectWeight:         r.EffectWeight,
        ServiceSource:        r.InstanceType,
        Scheme:               r.Scheme,
        Address:              r.Address,
        Port:                 r.Port,
        HealthCheckUrl:       r.HealthCheckUrl,
        Status:               r.Status,
        LastUpdatedTimestamp: r.LastUpdatedTimestamp,
        OverriddenStatus:     r.OverriddenStatus,
    }
    ins.init()
    return ins
}

type SQLRouteSource struct {
    KeyValueRouteSource
    driverName string
    url        string
    engine     *xorm.Engine
    routes     []*RouteModel
}

func NewSQLRouteSource(conf kvs.ConfigSource) *SQLRouteSource {
    driverName, err := conf.Get(SQL_DRIVERNAME) // "mysql"
    utils.Panic(err)

    url, err := conf.Get(SQL_URL) // root:kry02Local@?DB@tcp(172.16.1.248:3306)/po?charset=utf8
    utils.Panic(err)

    s := &SQLRouteSource{}
    s.driverName = driverName
    s.url = url
    s.conf = conf
    return s
}

func (p *SQLRouteSource) Init() {
    if p.isInit {
        return
    }
    engine, err := xorm.NewEngine(p.driverName, p.url)
    utils.Panic(err)
    p.engine = engine
    err = engine.Sync2(new(RouteModel))
    utils.Panic(err)
    if p.conf.GetBoolDefault(KEY_SQL_DISCOVERY_ENABLED, false) {
        err = engine.Sync2(new(HostInstanceModel))
        utils.Panic(err)
    }

    p.initRoutes()
    p.Watching(10 * time.Second)
    p.isInit = true
}

func (h *SQLRouteSource) initRoutes() {
    h.GetServicesInTime()
}

func (h *SQLRouteSource) CheckHealth(rootHealth *health.RootHealth) {
    if h.engine != nil {
        err := h.engine.Ping()
        if err == nil {
            rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_UP, Desc: "OK"}
        } else {
            rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_DOWN, Desc: err.Error()}
            rootHealth.Status = health.STATUS_DOWN
        }
    }
}

func (d *SQLRouteSource) Watching(second time.Duration) {
    d.run()
    go d.runTask(second)
}

func (d *SQLRouteSource) runTask(second time.Duration) {
    timer := time.NewTicker(second)
    for {
        select {
        case <-timer.C:
            go d.run()
        }
    }
}
func (d *SQLRouteSource) run() {
    routes, err := d.GetServicesInTime()
    if err == nil || routes != nil {
        d.execCallbacks(routes)
    } else {
        log.Error(err)
    }
}

func (d *SQLRouteSource) execCallbacks(routes []*RouteModel) {
    if d.routeChangedCallback != nil {
        for _, route := range routes {
            go d.routeChangedCallback(route.toRoute())
        }

    }
}

func (h *SQLRouteSource) GetServicesInTime() ([]*RouteModel, error) {
    var routes []*RouteModel
    //select * from routes where enabled=1 order by sorted desc limit 0,100
    err := h.engine.Where("enabled=?", 1).Desc("sorted").Limit(100, 0).Find(&routes)
    if err != nil {
        return nil, err
    }
    sets := mapset.NewSet()

    for _, route := range routes {
        h.Add(route.toRoute())
        log.WithField("route", route).Debug("add router by sql: ")
        sets.Add(route.ServiceId)
    }
    if h.conf.GetBoolDefault(KEY_SQL_DISCOVERY_ENABLED, false) {
        for _, serviceId := range sets.ToSlice() {
            hises := h.findHostInstances(serviceId.(string))
            //for _, ins := range hises {
            //     h.hostChangedCallback(serviceId.(string), ins.toHostInstance())
            //}

            size := len(hises)
            hosts := make([]*HostInstance, size)
            for i, ins := range hises {
                //h.hostChangedCallback(appName, newHostInstanceByEureka(appName, &instance))
                hosts[i] = ins.toHostInstance()
            }
            h.appHostsChangedCallback(serviceId.(string), hosts)

        }
    }

    return routes, nil

}

func (h *SQLRouteSource) findHostInstances(serviceId string) []*HostInstanceModel {
    hises := make([]*HostInstanceModel, 0)
    err := h.engine.Where("enabled=?", 1).And("service_id=?", serviceId).Find(&hises)
    if err != nil {
        return nil
    }
    return hises
}

func (h *SQLRouteSource) findRoutes() []*RouteModel {
    total := make([]*RouteModel, 0)
    limit, offset := 100, 0
    for {
        var routes []*RouteModel
        //select * from routes where enabled=1 order by sorted desc limit 0,100
        err := h.engine.Where("enabled=?", 1).Desc("sorted").Limit(limit, offset).Find(&routes)
        if err != nil {
            continue
        }
        if len(routes) < limit {
            break
        }
        offset += limit
        total = append(total, routes...)
    }
    return total

}

func (c *SQLRouteSource) GetServices() ([]*RouteModel) {
    if c.routes == nil {
        routes, err := c.GetServicesInTime()
        if err == nil {
            return routes
        }
    }
    return c.routes
}
