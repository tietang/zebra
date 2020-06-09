package proxy

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
	Method  string
	Path    string
	Handler Handler
}

type HttpServerRouter struct {
	Endpoints   []Endpoint
	ContextPath string
}

func (h *HttpServerRouter) Add(method, pathStr string, handler Handler) {
	e := Endpoint{
		Id:      strings.Join([]string{method, pathStr}, ":"),
		Method:  method,
		Path:    path.Join(h.ContextPath, pathStr),
		Handler: handler,
	}
	//Add(h, e)

	hasExists := false
	for i, endpoint := range h.Endpoints {
		if endpoint.Id == e.Id {
			log.Info("endpoint updated:", i, ": ", endpoint.Id, " = ", endpoint.Id)
			h.Endpoints[i] = e
			//rts := append(h.Endpoints[:i], e)
			//rts = append(rts, h.Endpoints[i+1:]...)
			//h.Endpoints = rts
			hasExists = true
		}
	}
	if len(h.Endpoints) == 0 || !hasExists {
		index := len(h.Endpoints)
		log.Info("endpoint add:", index, ": ", e.Id)
		h.Endpoints = append(h.Endpoints, e)
	}

}

func (h *HttpServerRouter) endpointExec(method, path string, ctx *Context) bool {
	e := h.find(method, path)
	if e == nil {
		return false
	}
	err := e.Handler(ctx)
	if err != nil {
		ctx.SetStatusCode(http.StatusInternalServerError)
		ctx.ResponseWriter.Write([]byte(err.Error()))
	}
	return true
}

func (h *HttpServerRouter) find(method, path string) *Endpoint {
	for _, endpoint := range h.Endpoints {
		isMethodMatched := strings.ToUpper(endpoint.Method) == strings.ToUpper(method) || strings.ToUpper(endpoint.Method) == MethodAny
		if isMethodMatched && endpoint.Path == path {
			return &endpoint
		}
	}
	return nil
}

func (h *HttpServerRouter) Any(path string, handler Handler) {
	h.Add(MethodAny, path, handler)
}
func (h *HttpServerRouter) Get(path string, handler Handler) {
	h.Add(MethodGet, path, handler)
}
func (h *HttpServerRouter) Post(path string, handler Handler) {
	h.Add(MethodPost, path, handler)
}
func (h *HttpServerRouter) Put(path string, handler Handler) {
	h.Add(MethodPut, path, handler)
}
func (h *HttpServerRouter) Delete(path string, handler Handler) {
	h.Add(MethodDelete, path, handler)
}
func (h *HttpServerRouter) Connect(path string, handler Handler) {
	h.Add(MethodConnect, path, handler)
}
func (h *HttpServerRouter) Head(path string, handler Handler) {
	h.Add(MethodHead, path, handler)
}
func (h *HttpServerRouter) Patch(path string, handler Handler) {
	h.Add(MethodPatch, path, handler)
}
func (h *HttpServerRouter) Options(path string, handler Handler) {
	h.Add(MethodOptions, path, handler)
}
func (h *HttpServerRouter) Trace(path string, handler Handler) {
	h.Add(MethodTrace, path, handler)
}
