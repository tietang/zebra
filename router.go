package zuul

import (
    "strings"
)

type Router struct {
    StripPrefix bool
    Routes      []Route
}

type Route struct {
    App                string
    Source             string
    SourcePrefix       string
    SourceIsFuzzyMatch bool
    Target             string
    TargetIsFuzzyMatch bool
    TargetPrefix       string
    StripPrefix        bool
}

//	r := &Router{}
//	r.AddRoute(Route{Source: "/app1/v1/user", App: "app1", Target: "/v1/user", StripPrefix: false})
//	r.AddRoute(Route{Source: "/app1/v2/user", App: "app1", Target: "/app1/v2/user", StripPrefix: false})
//	r.AddRoute(Route{Source: "/app2/**", App: "app2", Target: "/app20/**", StripPrefix: true})
//	r.AddRoute(Route{Source: "/app3/**", App: "app3", StripPrefix: false})
//	r.AddRoute(Route{Source: "/app4/**", App: "app4", StripPrefix: true})

func (r *Router) AddRoute(route Route) {
    initRoute(&route)
    r.Routes = append(r.Routes, route)
    //fmt.Println(route.App, route.Source, route.Target)
}

func initRoute(route *Route) {
    sourceFuzzyMatchIndex := strings.Index(route.Source, "/**")
    if sourceFuzzyMatchIndex > -1 {
        route.SourcePrefix = strings.TrimSuffix(route.Source, "/**")
        route.SourceIsFuzzyMatch = true
    }
    targetFuzzyMatchIndex := strings.Index(route.Target, "/**")
    if targetFuzzyMatchIndex > -1 {
        route.TargetPrefix = strings.TrimSuffix(route.Target, "/**")
        route.TargetIsFuzzyMatch = true
    }
}

func (r *Router) SetRoutes(routes []Route) {
    r.Routes = routes
    for _, v := range r.Routes {
        initRoute(&v)
    }
}

func (r *Router) GetMatchRouteTargetPath(sourcePath string) string {
    route := r.GetMatchRoute(sourcePath)
    return getRouteTargetPath(route, sourcePath)
}

func (r *Router) GetMatchRoute(sourcePath string) *Route {
    for _, v := range r.Routes {
        if isMatch(sourcePath, &v) {
            return &v
        }
    }
    return nil
}

func getRouteTargetPath(route *Route, sourcePath string) string {

    tpath := sourcePath
    isStrip := false

    if route.Target != "" {

        if route.TargetIsFuzzyMatch {
            tpath = sourcePath
            isStrip = true
        } else {
            return route.Target
        }
    } else {
        isStrip = true
    }

    if isStrip && route.StripPrefix {

        //print(strings.format("%s %d %d",tpath,index,strings.len(route.prefix)))
        if route.TargetIsFuzzyMatch {
            return route.TargetPrefix + strings.TrimPrefix(tpath, route.SourcePrefix)
        } else {
            return strings.TrimPrefix(tpath, route.SourcePrefix)
        }
    } else {
        if route.TargetIsFuzzyMatch {
            return route.TargetPrefix + tpath
        } else {
            return tpath
        }
    }
}

func isMatch(sourcePath string, route *Route) bool {

    if route.SourceIsFuzzyMatch {
        hasPrefix := strings.HasPrefix(sourcePath, route.SourcePrefix)
        if hasPrefix {
            return true
        }
    } else {
        if sourcePath == route.Source {
            return true
        }
    }
    return false
}
