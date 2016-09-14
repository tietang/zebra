package zuul

import (
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "strconv"
    "strings"
    "time"

    "github.com/tietang/go-eureka-client/eureka"
    "github.com/tietang/hystrix-go/hystrix"
    "github.com/tietang/stats"
)

type Handler struct {
    Counter              *stats.Counter
    Router               *Router
    DiscoveryRobin       *DiscoveryRobin
    httpClient           *http.Client
    IsCircuit            bool
    hystrixStreamHandler *hystrix.StreamHandler
    //config
    Timeout              time.Duration
    Threshold            int64
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {

    //	time.Sleep(time.Millisecond * 100)
    url := req.URL
    if strings.Contains(url.Path, "/hystrix.stream") {
        h.hystrixStreamHandler.ServeHTTP(w, req)
        return
    }
    //	w.Write([]byte(url.Path + "  |  " + url.RawPath + "  |  " + url.RawQuery))
    fmt.Println(url.Path + "  |  " + url.RawPath + "  |  " + url.RawQuery)
    if h.httpClient == nil {
        h.httpClient = &http.Client{
            Timeout: h.Timeout,
        }
    }

    if h.IsCircuit {
        h.hystrix(w, req)
    } else {
        h.forward(w, req)
    }

}

func (h *Handler) hystrix(w http.ResponseWriter, request *http.Request) {
    _, appName, targetPath, err := h.locate(w, request)
    method := request.Method
    url := request.URL
    path := url.Path
    if h.handleError(err, path, w) {
        return
    }

    cmdKey := appName + ":" + method + ":" + targetPath
    cmdKey = appName
    hystrix.Do(cmdKey, func() error {
        // talk to other services
        h.forward(w, request)
        return nil
    }, func(err error) error {
        msg := " fallback error:" + err.Error()
        message := NewMessage(http.StatusInternalServerError, msg, targetPath, err)
        // do this when services are down
        body, _ := json.Marshal(message)
        w.Write(body)
        return nil // newFallbackError(errFallback, msg)
    })

}

func (h *Handler) locate(w http.ResponseWriter, request *http.Request) (string, string, string, error) {
    url := request.URL
    path := url.Path
    queryStr := url.RawQuery
    appName, targetPath, err := h.route(path)
    return queryStr, appName, targetPath, err
}

func (h *Handler) handleError(err error, targetPath string, w http.ResponseWriter) bool {
    if err != nil {
        msg := " route error:" + err.Error()
        message := NewMessage(http.StatusInternalServerError, msg, targetPath, err)
        // do this when services are down
        body, _ := json.Marshal(message)
        w.Write(body)
        return true
    }

    return false
}

func (h *Handler) forward(w http.ResponseWriter, request *http.Request) {
    //	url := request.URL
    //	path := url.Path
    //	queryStr := url.RawQuery
    //	appName, targetPath, err := h.route(path)
    queryStr, appName, targetPath, err := h.locate(w, request)
    if h.handleError(err, targetPath, w) {
        return
    }
    _, _, ins := h.DiscoveryRobin.Next(appName)
    if ins == nil {
        w.Write([]byte("not found available instance for " + appName))
        return
    }

    //
    urlStr := h.ins2url(ins, targetPath, queryStr)

    //定义请求
    req, err := http.NewRequest("GET", urlStr, nil)
    //复制header
    for k, values := range request.Header {
        for _, v := range values {
            req.Header.Add(k, v)
        }
    }
    req.Header.Add("X-Forward-eureka", "")
    //调用请求
    res, err := h.httpClient.Do(req)
    if err != nil {
        fmt.Println(err)
        return
    }
    //处理response
    respBody, err := ioutil.ReadAll(res.Body)
    w.WriteHeader(res.StatusCode)
    for k, values := range res.Header {
        for _, v := range values {
            w.Header().Add(k, v)
        }
    }
    if err != nil {
        fmt.Println(err)
        return
    }

    size, err := w.Write(respBody)

    if size == len(respBody) {

    }

}

func (h *Handler) ins2url(ins *eureka.InstanceInfo, targetPath string, queryStr string) string {

    var port int
    schema := "http"
    if ins.Port.Enabled {
        port = ins.Port.Port
    } else {
        port = ins.SecurePort.Port
        schema = "https"
    }
    //	if strings.() == "" {

    //	}
    url := schema + "://" + ins.IpAddr + ":" + strconv.Itoa(port) + targetPath + "?" + queryStr
    return url

}

func (h *Handler) route(path string) (string, string, error) {

    route := h.Router.GetMatchRoute(path)
    if route == nil {
        return "", "", newNotFoundRouteError(errNotFoundRouteError, "not found route for path: " + path)
    } else {
        //		fmt.Println(route)
    }

    targetPath := getRouteTargetPath(route, path)
    appName := route.App
    return appName, targetPath, nil
}
