package router

import (
    log "github.com/sirupsen/logrus"
    "github.com/tietang/zebra/health"
    "reflect"
)

const ()

type Router struct {
    Name        string
    StripPrefix bool
    Routes      []Route
    //middleware  Handlers
    routeSources   []RouteSource
    healthCheckers []health.HealthChecker
}

func NewRouter() *Router {

    r := &Router{
        Name:         "default",
        Routes:       make([]Route, 0),
        routeSources: make([]RouteSource, 0),
    }
    return r
}

func (r *Router) register(routeSource RouteSource) {
    if routeSource == nil {
        return
    }
    routeSource.SetRouterChangedCallback(func(route *Route) {
        r.AddRoute(*route)
    })
    r.routeSources = append(r.routeSources, routeSource)
    r.healthCheckers = append(r.healthCheckers, routeSource.(health.HealthChecker))
    log.Info("RouteSource register: ", routeSource.Name())
}

func (r *Router) Start() {

    for _, routeSource := range r.routeSources {
        go func(rs RouteSource) {
            rs.Init()
            name := reflect.TypeOf(rs).Elem().Name()
            log.Info(name, ": ", " is starting...[", rs.Name(), "]")
            rs.Start()
            log.Info(name, " is started")
        }(routeSource)
    }
}

//	r := &Router{}
//	r.AddRoute(Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false})
//	r.AddRoute(Route{Source: "/app1/v2/user", ServiceId: "app1", Target: "/app1/v2/user", StripPrefix: false})
//	r.AddRoute(Route{Source: "/app2/**", ServiceId: "app2", Target: "/app20/**", StripPrefix: true})
//	r.AddRoute(Route{Source: "/app3/**", ServiceId: "app3", StripPrefix: false})`
//	r.AddRoute(Route{Source: "/app4/**", ServiceId: "app4", StripPrefix: true})

func (r *Router) AddRoute(route Route) {
    route.Init()
    hasExists := false

    for i, rt := range r.Routes {

        if rt.Id == route.Id {
            log.Debug("route updated: ", i, ": ", rt.Id, " = ", route.Id)
            //rts := append(r.Routes[:i], route)
            //rts = append(rts, r.Routes[i+1:]...)
            //r.Routes = rts
            r.Routes[i] = route
            hasExists = true
        }
    }
    if len(r.Routes) == 0 || !hasExists {
        log.Debug("route append: ", route.Id)
        r.Routes = append(r.Routes, route)
    }

    //fmt.Println(len(r.Routes), route.ServiceId, route.Source, route.Target)
}

func (r *Router) SetRoutes(routes []Route) {
    r.Routes = routes
    for _, v := range r.Routes {
        v.Init()
    }
}

func (r *Router) GetMatchRouteTargetPath(sourcePath string) string {
    route := r.GetMatchRoute(sourcePath)
    return route.GetRouteTargetPath(sourcePath)
}

func (r *Router) GetMatchRoute(sourcePath string) *Route {
    for _, v := range r.Routes {
        if v.isMatch(sourcePath) {
            return &v
        }
    }
    return nil
}
func (r *Router) GetHealthCheckers() []health.HealthChecker {
    return r.healthCheckers
}

//
//func (r *Router) CheckHealth(rootHealth *health.RootHealth) {
//    for _, routeSource := range r.routeSources {
//        routeSource.CheckHealth(rootHealth)
//    }
//}

// Use appends Handler(s) to the current Party's routes and child routes.
// If the current Party is the root, then it registers the middleware to all child Parties' routes too.
//func (r *Router) Use(handlers ...Handler) {
//    r.middleware = append(r.middleware, handlers...)
//}

// UseGlobal registers Handler middleware  to the beginning, prepends them instead of append
//
// Use it when you want to add a global middleware to all parties, to all routes in  all subdomains
// It should be called right before Listen functions
//func (r *Router) UseGlobal(handlers ...Handler) {
//    for _, route := range r.Routes {
//        route.Handlers = append(handlers, route.Handlers...) // prepend the handlers
//    }
//    r.middleware = append(handlers, r.middleware...) // set as middleware on the next routes too
//}
