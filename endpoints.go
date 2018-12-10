package zebra

import (
    log "github.com/sirupsen/logrus"
    "net/http"
    "path"
    "strings"
)

// These may or may not stay, you can use net/http's constants too.
const (
    MethodGet     = "GET"
    MethodPost    = "POST"
    MethodPut     = "PUT"
    MethodDelete  = "DELETE"
    MethodConnect = "CONNECT"
    MethodHead    = "HEAD"
    MethodPatch   = "PATCH"
    MethodOptions = "OPTIONS"
    MethodTrace   = "TRACE"
    //
    MethodAny = "ANY"
    // MethodNone is declared at iris.go, it will stay.
)

type Endpoint struct {
    Id      string
    method  string
    path    string
    handler Handler
}

func (h *HttpProxyServer) Add(method, pathStr string, handler Handler) {
    e := Endpoint{
        Id:      strings.Join([]string{method, pathStr}, ":"),
        method:  method,
        path:    path.Join(h.contextPath, pathStr),
        handler: handler,
    }
    Add(h, e)

}

func (h *HttpProxyServer) endpointExec(method, path string, ctx *Context) bool {
    e := h.find(method, path)
    if e == nil {
        return false
    }
    err := e.handler(ctx)
    if err != nil {
        ctx.SetStatusCode(http.StatusInternalServerError)
        ctx.ResponseWriter.Write([]byte(err.Error()))
    }
    return true
}

func (h *HttpProxyServer) find(method, path string) *Endpoint {
    for _, endpoint := range h.endpoints {
        isMethodMatched := strings.ToUpper(endpoint.method) == strings.ToUpper(method) || strings.ToUpper(endpoint.method) == MethodAny
        if isMethodMatched && endpoint.path == path {
            return &endpoint
        }
    }
    return nil
}

func Add(h *HttpProxyServer, e Endpoint) {
    hasExists := false
    for i, endpoint := range h.endpoints {
        if endpoint.Id == e.Id {
            log.Info("endpoint updated:", i, ": ", endpoint.Id, " = ", endpoint.Id)
            h.endpoints[i] = e
            //rts := append(h.endpoints[:i], e)
            //rts = append(rts, h.endpoints[i+1:]...)
            //h.endpoints = rts
            hasExists = true
        }
    }
    if len(h.endpoints) == 0 || !hasExists {
        index := len(h.endpoints)
        log.Info("endpoint add:", index, ": ", e.Id)
        h.endpoints = append(h.endpoints, e)
    }

    //fmt.Println(len(r.Routes), route.App, route.Source, route.Target)
}

func (h *HttpProxyServer) Any(path string, handler Handler) {
    h.Add(MethodAny, path, handler)
}
func (h *HttpProxyServer) Get(path string, handler Handler) {
    h.Add(MethodGet, path, handler)
}
func (h *HttpProxyServer) Post(path string, handler Handler) {
    h.Add(MethodPost, path, handler)
}
func (h *HttpProxyServer) Put(path string, handler Handler) {
    h.Add(MethodPut, path, handler)
}
func (h *HttpProxyServer) Delete(path string, handler Handler) {
    h.Add(MethodDelete, path, handler)
}
func (h *HttpProxyServer) Connect(path string, handler Handler) {
    h.Add(MethodConnect, path, handler)
}
func (h *HttpProxyServer) Head(path string, handler Handler) {
    h.Add(MethodHead, path, handler)
}
func (h *HttpProxyServer) Patch(path string, handler Handler) {
    h.Add(MethodPatch, path, handler)
}
func (h *HttpProxyServer) Options(path string, handler Handler) {
    h.Add(MethodOptions, path, handler)
}
func (h *HttpProxyServer) Trace(path string, handler Handler) {
    h.Add(MethodTrace, path, handler)
}
