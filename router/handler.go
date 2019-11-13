package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/rcrowley/go-metrics"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/hystrix-go/hystrix"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/httpclient"
	"github.com/tietang/zebra/utils"
	"io"
	"math"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"sync"
	"time"
)

const (
	METER_ERROR_REQUEST_PREFIX = "error:"
	METER_OK_REQUEST_PREFIX    = "ok:"
	METER_50X_REQUEST_PREFIX   = "50x:"
	METER_40X_REQUEST_PREFIX   = "40x:"
	KEY_SERVER_MODE            = "server.mode"
	MODE_CLIENT                = "client"
	MODE_REVERSE_PROXY         = "reverseproxy"
	DEFAULT_MODE               = MODE_CLIENT
	SERVER_DEBUG               = "server.debug"
	HTTP_PROXY_SERVER_NAME     = "zebra"
)

var DefaultRouters []*Router
var DefaultHosts *sync.Map
var DefaultUnavailableHosts *sync.Map

type RequestContext struct {
	Writer       http.ResponseWriter
	Request      *http.Request
	QueryStr     string
	AppName      string
	TargetPath   string
	HostInstance *HostInstance
	Error        error
	//
	rawUrl           string
	statusCode       int
	responseByteSize int64

	logObjects map[string]interface{}

	//
	isCanRetry    bool
	isGrayRequest bool
}

type UniversalHandler struct {
	HttpClients      *httpclient.HttpClients
	conf             kvs.ConfigSource
	Balancer         *DiscoveryBalancer
	Router           *Router
	RequestCondition RequestCondition
}

func NewUniversalHandler(conf kvs.ConfigSource) *UniversalHandler {
	rh := &UniversalHandler{
		conf: conf,
	}
	//初始化meter
	rh.HttpClients = httpclient.NewHttpClients(conf)
	rh.Balancer = NewDiscoveryBalancer(conf)
	rh.Router = NewRouter()
	DefaultRouters = append(DefaultRouters, rh.Router)
	DefaultHosts = rh.Balancer.Hosts
	DefaultUnavailableHosts = rh.Balancer.UnavailableHosts
	//
	if conf.GetBoolDefault(KEY_EUREKA_ENABLED, false) {
		log.Info("eureka discovery enabled.")
		ers := NewEurekaRouteSource(conf)
		rh.Router.register(ers)
		rh.Balancer.register(ers)
	}
	//
	//if conf.GetBoolDefault(KEY_K8S_ENABLED, false) {
	//	log.Info("k8s discovery enabled.")
	//	prs := NewKubernetesRouteSource(conf)
	//	rh.Router.register(prs)
	//	rh.Balancer.register(prs)
	//}
	//
	//if conf.GetBoolDefault(KEY_CONSUL_ENABLED, false) {
	//	log.Info("consul discovery&routes enabled.")
	//	crs := NewConsulRouteSource(conf)
	//	rh.Router.register(crs)
	//	rh.Balancer.register(crs)
	//}

	if conf.GetBoolDefault(KEY_ZK_ENABLED, false) {
		log.Info("zookeeper routes enabled.")
		zrs := NewZookeeperRouteSource(conf)
		rh.Router.register(zrs)
		rh.Balancer.register(zrs)
	}
	//
	if conf.GetBoolDefault(KEY_INI_ROUTES_ENABLED, false) {
		log.Info("ini file routes enabled.")

		prs := NewIniFileRouteSource(conf)
		rh.Router.register(prs)
		rh.Balancer.register(prs)
	}

	if conf.GetBoolDefault(KEY_SQL_ROUTES_ENABLED, false) {
		log.Info("sql discovery&routes enabled.")
		prs := NewSQLRouteSource(conf)
		rh.Router.register(prs)
		rh.Balancer.register(prs)
	}

	if conf.GetBoolDefault("traffic.cond.enabled", false) {
		typ, err := conf.Get("traffic.cond.type")
		if err != nil || typ == "" {
			log.Info("not config traffic.cond.type")
		} else {
			if typ == "composite" {
				typs := conf.Strings("traffic.cond.composites")
				c := new(CompositeRequestCondition)
				for _, typ := range typs {
					rc := RequestConditions[typ]
					c.Add(rc)
				}
				rh.RequestCondition = c
			} else {
				rc := RequestConditions[typ]
				rh.RequestCondition = rc
			}
			rh.RequestCondition.Conf(conf)
		}

	}

	//外部注册扩展
	log.Info("register ext & plugins: ")
	for i, rs := range globalRouterSources {
		rs.Build()
		rh.Router.register(rs)
		rh.Balancer.register(rs)
		log.Info(i, ": ", rs.Name())
	}
	rh.Router.Start()
	return rh
}

func (h *UniversalHandler) IsCircuit() bool {
	return h.conf.GetBoolDefault(KEY_CIRCUIT_ENABLED, true)
}

func (h *UniversalHandler) Handle(w http.ResponseWriter, req *http.Request) bool {
	url := req.URL
	log.Debug(url.Path, "  |  ", url.RawPath, "  |  ", url.RawQuery)

	s := time.Now()
	//在HttpProxyServer.ServeHTTP中已经包含
	//urlTimer := metrics.GetOrRegisterTimer(url.Path, UrlRegistry)
	ctx := h.locate(w, req)
	serviceTimer := metrics.GetOrRegisterTimer(ctx.AppName, ServiceRegistry)
	//ctx := &RequestContext{
	//    Writer:     w,
	//    Request:    req,
	//    QueryStr:   queryStr,
	//    AppName:    appName,
	//    TargetPath: targetPath,
	//    Error:      err,
	//}
	if h.handleError("route error,", ctx.Error, ctx) {
		return false
	}
	//if ctx.Error != nil {
	//    return false
	//}
	serviceTimer.Time(func() {
		//urlTimer.Time(func() {
		if h.IsCircuit() {
			err := h.hystrix(ctx)
			if err != nil {
				log.Error(err)
			}
		} else {
			err := h.forward(ctx)
			if err != nil {
				log.Error(err)
			}
		}
		//})
	})

	nano := time.Since(s)
	ms := nano.Nanoseconds() / int64(time.Millisecond)
	mcs := nano.Nanoseconds() / int64(time.Microsecond)
	log.Debug("res nano time:", ms, "    ", mcs)
	//d := time.Duration(nano) * time.Nanosecond
	//serviceTimer.Update(d)
	//urlTimer.Update(d)
	//$remote_addr - $remote_user [$time_local] '
	//                   '"$request" $status $bytes_sent '
	//'"$http_referer" "$http_user_agent" "$gzip_ratio
	//127.0.0.1 - - [05/Dec/2016:14:41:29 +0800] "GET /_admin/qps.json HTTP/1.1" 200 99 "http://127.0.0.1:8000/_dashboard/chart/area2.html" "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_12_1) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/54.0.2840.98 Safari/537.36"
	if ctx.Error == nil {
		log.WithFields(ctx.logObjects).Info(ctx.statusCode, ms, "MS ", ctx.AppName, " ", ctx.responseByteSize, "b ", ctx.rawUrl)
	} else {
		//2016/09/26 16:59:36 [error] 21008#0: *4568782 kevent() reported that connect() failed (61: Connection refused) while connecting to upstream, client: 127.0.0.1, server: localhost, request: "GET /api/health HTTP/1.1", upstream: "http://192.168.4.198:7912/health", host: "127.0.0.1:8000"
		log.WithFields(ctx.logObjects).Error(ctx.statusCode, " ", ctx.AppName, " ", ctx.Error)
	}
	return true
}

func (h *UniversalHandler) hystrix(ctx *RequestContext) error {
	req := ctx.Request
	_, appName, targetPath, _ := ctx.QueryStr, ctx.AppName, ctx.TargetPath, ctx.Error

	method := req.Method
	//if h.handleError("route error: ", err, ctx) {
	//    return err
	//}
	cmdKey := appName + ":" + method + ":" + targetPath
	cmdKey = appName
	//Synchronous invoke
	error := hystrix.Do(cmdKey, func() error {
		// talk to other services
		return h.innerForward(ctx)
	}, func(err error) error {
		h.handleError("fallback: ", err, ctx)
		return nil //newFallbackError(errFallback, msg)
	})
	if error != nil {
		log.Error("hystrix exec: ", error)
		return error
	}
	return nil

}

func (h *UniversalHandler) handleError(prefixMsg string, err error, ctx *RequestContext) bool {
	if err != nil {
		ctx.Error = err
		errorStr := err.Error()

		statusCode := http.StatusInternalServerError
		msgSuffix := ""
		//IsAvailable := false
		if strings.Contains(errorStr, "timeout") {
			statusCode = http.StatusGatewayTimeout
			msgSuffix = " timeout."
		} else if strings.Contains(errorStr, "i/o timeout") { //
			msgSuffix = " connect timeout."
			statusCode = http.StatusServiceUnavailable
			//IsAvailable = true
		} else if strings.Contains(errorStr, "connection refused") {
			msgSuffix = " connection refused."
			statusCode = http.StatusServiceUnavailable
			//IsAvailable = true
		} else {
			msgSuffix = " Internal GRPCServer Error."
			statusCode = http.StatusInternalServerError
		}
		msg := prefixMsg + " " + msgSuffix
		//log.Error(statusCode, "  ", msg, "  : ", ctx.TargetPath, "   ", err)

		message := utils.NewMessage(statusCode, msg, ctx.TargetPath, err)
		//
		if ctx.HostInstance != nil {

			//if IsAvailable {
			//    h.Balancer.AddUnavailableHostInstance(ctx.HostInstance.AppName, ctx.HostInstance)
			//}
			d := getAppFailTimeWindowSeconds(h.conf, ctx.HostInstance.AppName)
			m := GetOrRegisterErrorMeter(ctx.HostInstance, d)
			m.Mark(1)
		}

		if ctx.statusCode == 0 {
			ctx.statusCode = statusCode
		}
		// do this when services are down
		h.handlerStatusCode(statusCode, ctx)
		body, _ := json.Marshal(message)
		ctx.Writer.WriteHeader(statusCode)
		ctx.Writer.Write(body)

		return true
	}

	return false
}

func (h *UniversalHandler) innerForward(ctx *RequestContext) error {
	//负载均衡
	//_, hostIns, ins := h.DiscoveryRobin.Next(appName)\\
	path := ctx.Request.URL.Path
	isEnabledGray := h.conf.GetBoolDefault(fmt.Sprintf("traffic.cond.%s.enabled", ctx.AppName), false)
	if isEnabledGray && h.RequestCondition != nil && h.RequestCondition.Matched(ctx) {
		ctx.isGrayRequest = true
		ctx.HostInstance = h.Balancer.Next(ctx.AppName, path, true)
	} else {
		ctx.isGrayRequest = false
		ctx.HostInstance = h.Balancer.Next(ctx.AppName, path, false)
	}

	//var ins *eureka.InstanceInfo
	//ins = hostIns.ExternalInstance.(*eureka.InstanceInfo)

	if ctx.HostInstance == nil {
		msg := "not found available instance for " + ctx.AppName
		return errors.New(msg)
	}

	var err error
	instanceTimer := metrics.GetOrRegisterTimer(ctx.HostInstance.InstanceId, InstanceRegistry)
	instanceTimer.Time(func() {
		err = h.innerForward0(ctx)
	})

	//if err == nil {
	//    m := metrics.GetOrRegisterMeter(METER_OK_REQUEST_PREFIX+ctx.HostInstance.InstanceId, InstanceRegistry)
	//    m.Mark(1)
	//}
	return err
}
func (h *UniversalHandler) innerForward0(ctx *RequestContext) error {
	mode := h.conf.GetDefault(KEY_SERVER_MODE, DEFAULT_MODE)
	if strings.ToLower(mode) == MODE_CLIENT {
		return h.innerForwardClient(ctx)
	}

	if strings.ToLower(mode) == MODE_REVERSE_PROXY {
		return h.innerForwardReverseProxy(ctx)
	}
	//应该不会被执行
	return h.innerForwardClient(ctx)

}

func (h *UniversalHandler) innerForwardReverseProxy(ctx *RequestContext) error {

	//ins := ctx.HostInstance
	//TODO 如果proxy需要动态传参数，就需要构造targetQuery
	//targetQuery := ""
	//director := func(req *http.Request) {
	//    //req = ctx.Request
	//    req.URL.Scheme = ins.Scheme
	//    req.URL.Host = ins.Address + ":" + ins.Port
	//
	//    req.URL.Path = ctx.TargetPath // singleJoiningSlash(target.Path, req.URL.Path)
	//    if targetQuery == "" || req.URL.RawQuery == "" {
	//        req.URL.RawQuery = targetQuery + req.URL.RawQuery
	//    } else {
	//        req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	//    }
	//    if _, ok := req.Header["User-Agent"]; !ok {
	//        // explicitly disable User-Agent so it's not set to default value
	//        req.Header.Set("User-Agent", "")
	//    }
	//}
	//proxy := &httputil.ReverseProxy{Director: director}
	proxy := &httputil.ReverseProxy{Director: h.director0(ctx)}
	//utils.NewSingleHostReverseProxy(remote)
	ctx.Writer.Header().Add("X-Proxy-GRPCServer", "zebra")
	if h.conf.GetBoolDefault(SERVER_DEBUG, false) {
		ctx.Writer.Header().Add("X-Forward-Host", strings.Join([]string{ctx.HostInstance.Address, ctx.HostInstance.Port}, ":"))
	}
	if ctx.isGrayRequest {
		ctx.Writer.Header().Add("X-Gray-version", ctx.HostInstance.Version)
	}

	proxy.ModifyResponse = h.modifyResponse(ctx)
	proxy.ServeHTTP(ctx.Writer, ctx.Request)
	//val := reflect.ValueOf(ctx.Writer)
	//if val.Elem().Type().String() == "http.response" {
	//    //fmt.Println(val.Elem().Type())
	//    valStatue := val.Elem().FieldByName("StatusCode")
	//    contentLength := val.Elem().FieldByName("ContentLength")
	//    //fmt.Println(valStatue.Int())
	//    statusCode := int(valStatue.Int())
	//
	//    h.handlerStatusCode(statusCode, ctx)
	//    ctx.responseByteSize = int64(contentLength.Int())
	//
	//}
	return nil
}
func (h *UniversalHandler) modifyResponse(ctx *RequestContext) func(res *http.Response) error {
	return func(res *http.Response) error {
		h.handlerStatusCode(res.StatusCode, ctx)
		ctx.responseByteSize = res.ContentLength
		return nil
	}
}

func (h *UniversalHandler) handlerStatusCode(statusCode int, ctx *RequestContext) {
	if statusCode >= 400 && statusCode < 500 {
		m := metrics.GetOrRegisterMeter(METER_40X_REQUEST_PREFIX+ctx.AppName, ServiceRegistry)
		m.Mark(1)
	}
	if statusCode >= 500 {
		m := metrics.GetOrRegisterMeter(METER_50X_REQUEST_PREFIX+ctx.AppName, ServiceRegistry)
		m.Mark(1)
		if ctx.HostInstance != nil {
			m := metrics.GetOrRegisterMeter(METER_ERROR_REQUEST_PREFIX+ctx.AppName, ServiceRegistry)
			m.Mark(1)
		}
	}
	ctx.statusCode = statusCode
}

func (h *UniversalHandler) director0(ctx *RequestContext) func(req *http.Request) {

	ins := ctx.HostInstance
	//TODO 如果proxy需要动态传参数，就需要构造targetQuery
	targetQuery := ""
	return func(req *http.Request) {
		req.URL.Scheme = ins.Scheme
		req.URL.Host = ins.Address + ":" + ins.Port

		req.URL.Path = ctx.TargetPath
		//req.URL.Path =  singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}
		req.Header.Add("X-Proxy-server", HTTP_PROXY_SERVER_NAME)
		if ctx.isGrayRequest {
			req.Header.Add("X-Gray", "true")
		}
		ctx.rawUrl = req.URL.String()
		log.Debug("backend url: ", req.URL.String())
	}
}

func (h *UniversalHandler) director1(ctx *RequestContext) func(req *http.Request) {
	urlStr := h.ins2url(ctx.HostInstance, ctx.TargetPath, ctx.QueryStr)
	ctx.rawUrl = urlStr
	remote, err := url.Parse(urlStr)
	if err != nil {
		h.handleError("remote url parse error", err, ctx)
	}
	log.Debug("remote: ", remote)
	target := remote
	targetQuery := target.RawQuery
	return func(req *http.Request) {
		req.URL.Scheme = target.Scheme
		req.URL.Host = target.Host
		req.URL.Path = target.Path

		//req.URL.Path =  singleJoiningSlash(target.Path, req.URL.Path)
		if targetQuery == "" || req.URL.RawQuery == "" {
			req.URL.RawQuery = targetQuery + req.URL.RawQuery
		} else {
			req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
		}
		if _, ok := req.Header["User-Agent"]; !ok {
			// explicitly disable User-Agent so it's not set to default value
			req.Header.Set("User-Agent", "")
		}

		req.Header.Add("X-Proxy-server", HTTP_PROXY_SERVER_NAME)
		if ctx.isGrayRequest {
			req.Header.Add("X-Gray", "true")
		}
		ctx.rawUrl = req.URL.String()
	}
}

func (h *UniversalHandler) innerForwardClient(ctx *RequestContext) error {
	//
	urlStr := h.ins2url(ctx.HostInstance, ctx.TargetPath, ctx.QueryStr)

	//targetQuery := ""
	//return ins.Scheme, ins.Address + ":" + ins.Port, targetPath, targetQuery
	method := ctx.Request.Method
	//定义请求
	req, err := http.NewRequest(method, urlStr, ctx.Request.Body)
	//TODO 如果proxy需要动态传参数，就需要构造targetQuery
	targetQuery := ""
	if targetQuery == "" || req.URL.RawQuery == "" {
		req.URL.RawQuery = targetQuery + req.URL.RawQuery
	} else {
		req.URL.RawQuery = targetQuery + "&" + req.URL.RawQuery
	}
	ctx.rawUrl = req.URL.String()
	if clientIP, _, err := net.SplitHostPort(req.RemoteAddr); err == nil {
		if prior, ok := req.Header["X-Forwarded-For"]; ok {
			clientIP = strings.Join(prior, ", ") + ", " + clientIP
		}
		req.Header.Set("X-Forwarded-For", clientIP)
	}
	//设置超时
	//复制header
	for k, values := range req.Header {
		for _, v := range values {
			req.Header.Add(k, v)
		}
	}
	req.Header.Add("X-Proxy-server", utils.HTTP_PROXY_SERVER_NAME)
	if ctx.isGrayRequest {
		req.Header.Add("X-Gray", "true")
	}

	log.Debug("remote: ", req)
	//调用请求
	res, err := h.HttpClients.Do(ctx.AppName, req)

	//m := metrics.NewMeter()
	//metrics.Register("quux", m)
	if err != nil {
		//log.Error(err)
		//if res.StatusCode >= 500 {

		//}
		return err
	}

	h.handlerStatusCode(res.StatusCode, ctx)

	for k, values := range res.Header {
		for _, v := range values {
			ctx.Writer.Header().Add(k, v)
		}
	}
	ctx.Writer.Header().Add("X-Proxy-GRPCServer", "zebra")
	if h.conf.GetBoolDefault(SERVER_DEBUG, false) {
		ctx.Writer.Header().Add("X-Forward-Host", strings.Join([]string{ctx.HostInstance.Address, ctx.HostInstance.Port}, ":"))
	}
	if ctx.isGrayRequest {
		ctx.Writer.Header().Add("X-Gray-version", ctx.HostInstance.Version)
	}
	ctx.Writer.WriteHeader(res.StatusCode)

	// 如果出错就不需要close，因此defer语句放在err处理逻辑后面
	// Clients must call resp.Body.Close when finished reading resp.Body -- from golang doc
	defer res.Body.Close()

	//处理response
	//respBody, err := ioutil.ReadAll(res.Body)
	size, _ := io.Copy(ctx.Writer, res.Body)
	// Reset resp.Body so it can be use again
	//res.Body = ioutil.NopCloser(bytes.NewBuffer(respBody))

	//
	if err := res.Body.Close(); err != nil {
		log.Error(err)
	}
	ctx.responseByteSize = size
	log.Debug("response body size: ", ctx.responseByteSize, ", content length: ", res.ContentLength)

	return err
}

func (h *UniversalHandler) forward(ctx *RequestContext) error {

	err := h.innerForward(ctx)

	if h.handleError("route error: ", err, ctx) {
		return nil
	}
	return err
}

func (h *UniversalHandler) ins2url(ins *HostInstance, targetPath string, queryStr string) string {
	schema := ins.Scheme
	port := ins.Port
	url := schema + "://" + ins.Address + ":" + port + targetPath + "?" + queryStr
	return url

}

func (h *UniversalHandler) hostIns2Url(ins *HostInstance, targetPath string) (string, string, string, string) {
	targetQuery := ""
	return ins.Scheme, ins.Address + ":" + ins.Port, targetPath, targetQuery

}
func (h *UniversalHandler) locate(w http.ResponseWriter, req *http.Request) *RequestContext {
	url := req.URL
	path := url.Path
	queryStr := url.RawQuery
	appName, targetPath, err := h.route(path)

	ctx := &RequestContext{
		Writer:     w,
		Request:    req,
		QueryStr:   queryStr,
		AppName:    appName,
		TargetPath: targetPath,
		Error:      err,
	}
	return ctx
}

func (h *UniversalHandler) route(path string) (string, string, error) {
	log.Debug(path)
	route := h.Router.GetMatchRoute(path)
	if route == nil {
		return "", "", utils.NewNotFoundRouteError(utils.ErrNotFoundRouteError, "not found route for path: "+path)
	} else {
		log.WithField("route", route).Debug()
	}

	targetPath := route.GetRouteTargetPath(path)
	appName := route.ServiceId
	return appName, targetPath, nil
}

func NewEWMA10() metrics.EWMA {
	return metrics.NewEWMA(1 - math.Exp(-5.0/60.0/0.5))
}
