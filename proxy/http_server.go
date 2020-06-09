package proxy

import (
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/health"
	"github.com/tietang/zebra/router"
	"net/http"
	"strconv"
	"time"
)

const (
	Version = "0.1"
	Name    = "Zebra"
)

//server
const (
	//key
	KEY_SERVER_DEBUG            = "server.debug"
	KEY_SERVER_PORT             = "app.server.port"
	KEY_SERVER_CONTEXT_PATH     = "server.contextPath"
	KEY_SERVER_MODE             = "server.mode"
	KEY_SERVER_FAVICON_ICO_PATH = "server.favicon.ico.Path"
	KEY_SERVER_GMS_ENABLED      = "server.gms.enabled"
	KEY_SERVER_GMS_DOMAIN       = "server.gms.domain"
	KEY_SERVER_GMS_PORT         = "server.gms.port"
	KEY_SERVER_GMS_DEBUG        = "server.gms.debug"
	//default value
	DEFAULT_SERVER_DEBUG            = "true"
	DEFAULT_SERVER_PORT             = 19001
	DEFAULT_SERVER_MODE             = "client"
	DEFAULT_SERVER_FAVICON_ICO_PATH = "favicon.ico"
	DEFAULT_SERVER_GMS_ENABLED      = "true"
	DEFAULT_SERVER_GMS_DOMAIN       = ""
	DEFAULT_SERVER_GMS_PORT         = 17980
	DEFAULT_SERVER_GMS_DEBUG        = "true"
)

type HttpProxyServer struct {
	HttpServerRouter
	dir        string
	configFile string
	port       int
	conf       kvs.ConfigSource
	health     *health.RootHealth
	middleware Handlers

	//
	Stats *Stats
}

func NewHttpProxyServer(conf kvs.ConfigSource) *HttpProxyServer {
	port := conf.GetIntDefault(KEY_SERVER_PORT, DEFAULT_SERVER_PORT)
	rootHealth := &health.RootHealth{}
	rootHealth.Status = health.STATUS_UP
	rootHealth.Desc = "Gateway zebra"
	rootHealth.Healths = make(map[string]*health.Health)
	contextPath := conf.GetDefault(KEY_SERVER_CONTEXT_PATH, "/")
	h := &HttpProxyServer{
		conf:   conf,
		health: rootHealth,
		port:   port,
	}
	h.ContextPath = contextPath
	h.Use(func(context *Context) error {
		return nil
	})
	h.initEndpoints()
	return h
}

func (h *HttpProxyServer) Run(addr string) {
	server := &http.Server{
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		Addr:         addr,
		Handler:      h,
	}

	go h.StartStatsServer()
	//	http.ListenAndServe(addr, h.proxyHandler)
	log.Info("http proxy server address: ", addr)
	server.ListenAndServe()

}

func (h *HttpProxyServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	ctx := NewContext(w, req, h.middleware)
	ctx.Handlers = h.middleware
	//	time.Sleep(time.Millisecond * 100)
	url := req.URL
	t := metrics.GetOrRegisterTimer(url.Path, router.UrlRegistry)
	t.Time(func() {
		if h.endpointExec(ctx.Method(), ctx.Path(), ctx) {
			return
		}
		ctx.Next()

	})

}
func (h *HttpProxyServer) DefaultRun() {
	if h.port == 0 {
		h.port = DEFAULT_SERVER_PORT
	}
	h.RunByPort(h.port)
}

func (h *HttpProxyServer) RunByPort(port int) {
	h.Run(":" + strconv.Itoa(port))
}

// Use appends Handler(s) to the current Party's routes and child routes.
// If the current Party is the root, then it registers the middleware to all child Parties' routes too.
func (r *HttpProxyServer) Use(handlers ...Handler) {
	r.middleware = append(r.middleware, handlers...)
}

//
//// UseGlobal registers Handler middleware  to the beginning, prepends them instead of append
////
//// Use it when you want to add a global middleware to all parties, to all routes in  all subdomains
//// It should be called right before Listen functions
//func (r *HttpProxyServer) UseGlobal(handlers ...Handler) {
//    //for _, host := range r.Routes {
//    //    host.Handlers = append(handlers, host.Handlers...) // prepend the handlers
//    //}
//    r.middleware = append(handlers, r.middleware...) // set as middleware on the next routes too
//}
