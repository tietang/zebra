package zuul

import (
    "encoding/json"
    "fmt"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/tietang/hystrix-go/hystrix"
    "github.com/tietang/stats"
)

type HttpProxyServer struct {
    router  *Router
    handler *Handler
}

func NewHttpProxyServer(eurekaUrl string, routes []string) *HttpProxyServer {
    counter := stats.NewCounter()
    router := &Router{}
    discovery := &Discovery{eurekaUrl: eurekaUrl}
    robin := NewDiscoveryRobin(&Robin{}, discovery)
    discovery.ScheduleAtFixedRate(10 * time.Second)
    var timeout time.Duration
    timeout = 20 * time.Second
    hystrixStreamHandler := hystrix.NewStreamHandler()
    hystrixStreamHandler.Start()
    handler := &Handler{
        Counter:              counter,
        Router:               router,
        DiscoveryRobin:       robin,
        IsCircuit:            true,
        hystrixStreamHandler: hystrixStreamHandler,
        Timeout:              timeout,
        Threshold:            10,
    }
    server := &HttpProxyServer{handler: handler, router: router}
    server.initRouter(routes)
    server.initDiscoveryRouter(discovery)

    body, err := json.Marshal(server.router.Routes)
    fmt.Println(string(body))
    fmt.Println(err)
    return server
}

func (h *HttpProxyServer) initRouter(routes []string) {
    for _, v := range routes {
        vs := strings.Split(v, ",")
        sp := false
        if len(vs) == 4 {
            sp, _ = strconv.ParseBool(vs[3])
        }
        route := Route{
            Source:      vs[0],
            App:         vs[1],
            Target:      vs[2],
            StripPrefix: sp,
        }
        h.router.AddRoute(route)
    }
}

func (h *HttpProxyServer) initDiscoveryRouter(discovery *Discovery) {
    apps := discovery.apps
    if apps == nil || apps.Applications == nil {
        return
    }
    for _, a := range apps.Applications {
        appName := strings.ToLower(a.Name)
        route := Route{
            Source:      "/" + appName + "/**",
            App:         a.Name,
            StripPrefix: true,
        }
        h.router.AddRoute(route)
    }
}

func (h *HttpProxyServer) Run(addr string) {
    server := &http.Server{
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        Addr:         addr,
        Handler:      h.handler,
    }
    //	http.ListenAndServe(addr, h.handler)

    server.ListenAndServe()
}

func (h *HttpProxyServer) DefaultRun() {
    h.Run(":8002")
}

func (h *HttpProxyServer) RunByPort(port int) {
    h.Run(":" + strconv.Itoa(port))
}
