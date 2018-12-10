package zebra

import (
    "fmt"

    "encoding/json"
    "net/http"
    "net/url"
)

type Handler func(*Context) error
type Handlers []Handler

// combineHandlers merges two lists of handlers into a new list.
func combineHandlers(h1 []Handler, h2 []Handler) []Handler {
    hh := make([]Handler, len(h1)+len(h2))
    copy(hh, h1)
    copy(hh[len(h1):], h2)
    return hh
}

// SerializeFunc serializes the given data of arbitrary type into a byte array.
type SerializeFunc func(data interface{}) ([]byte, error)

// Context represents the contextual data and environment while processing an incoming HTTP request.
type Context struct {
    ResponseWriter http.ResponseWriter
    Request        *http.Request

    Serialize SerializeFunc          // the function serializing the given data of arbitrary type into a byte array.
    data      map[string]interface{} // data items managed by Get and Set
    index     int                    // the index of the currently executing handler in handlers
    handlers  []Handler              // the handlers associated with the current route
    //
    attrs      map[string]interface{}
    attrKeys   []string // list of route parameter names
    attrValues []string // list of parameter values corresponding to pnames
}

func NewContext(w http.ResponseWriter, req *http.Request, handlers []Handler) *Context {
    return &Context{
        ResponseWriter: w,
        Request:        req,
        handlers:       handlers,
        Serialize: func(data interface{}) ([]byte, error) {
            return json.Marshal(data)
        },
    }
}

// Param returns the named parameter value that is found in the URL path matching the current route.
// If the named parameter cannot be found, an empty string will be returned.
func (c *Context) Attr(key string) string {
    for i, n := range c.attrKeys {
        if n == key {
            return c.attrValues[i]
        }
    }
    return ""
}

// GetCookie returns cookie's value by it's name
// returns empty string if nothing was found.
func (ctx *Context) Cookie(name string) string {
    cookie, err := ctx.Request.Cookie(name)
    if err != nil {
        return ""
    }
    return cookie.Value
}

// Get returns the named data item previously registered with the context by calling Set.
// If the named data item cannot be found, nil will be returned.
func (c *Context) Get(name string) interface{} {
    return c.data[name]
}

// Set stores the named data item in the context so that it can be retrieved later.
func (c *Context) Set(name string, value interface{}) {
    if c.data == nil {
        c.data = make(map[string]interface{})
    }
    c.data[name] = value
}

// Next calls the rest of the handlers associated with the current route.
// If any of these handlers returns an error, Next will return the error and skip the following handlers.
// Next is normally used when a handler needs to do some postprocessing after the rest of the handlers
// are executed.
func (c *Context) Next() error {
    c.index++
    for n := len(c.handlers); c.index < n; c.index++ {
        if err := c.handlers[c.index](c); err != nil {
            return err
        }
    }
    return nil
}

// Abort skips the rest of the handlers associated with the current route.
// Abort is normally used when a handler handles the request normally and wants to skip the rest of the handlers.
// If a handler wants to indicate an error condition, it should simply return the error without calling Abort.
func (c *Context) Abort() {
    c.index = len(c.handlers)
}

// WriteData writes the given data of arbitrary type to the response.
// The method calls the Serialize() method to convert the data into a byte array and then writes
// the byte array to the response.
func (c *Context) WriteData(data interface{}) (err error) {
    var bytes []byte
    if bytes, err = c.Serialize(data); err == nil {
        _, err = c.ResponseWriter.Write(bytes)
    }
    return
}
func (c *Context) Write(data []byte) (int, error) {
    return c.ResponseWriter.Write(data)
}

func (c *Context) WriteString(str string) (int, error) {
    return c.ResponseWriter.Write([]byte(str))
}

// Method returns the request.Method, the client's http method to the server.
func (ctx *Context) Method() string {
    return ctx.Request.Method
}

// Path returns the full request path,
// escaped if EnablePathEscape config field is true.
func (ctx *Context) Path() string {
    return ctx.Request.URL.Path
}

// SetStatusCode sets response status code.
func (ctx *Context) SetStatusCode(statusCode int) {
    ctx.ResponseWriter.WriteHeader(statusCode)
}

// SetContentType sets response Content-Type.
func (ctx *Context) SetContentType(contentType string) {
    ctx.ResponseWriter.Header().Add("Content-Type", contentType)
}

// RequestURI returns RequestURI.
//
// This uri is valid until returning from RequestHandler.
func (ctx *Context) RequestURI() string {
    return ctx.Request.RequestURI
}

// URI returns requested uri.
//
// The uri is valid until returning from RequestHandler.
func (ctx *Context) URL() *url.URL {
    return ctx.Request.URL
}

// Referer returns request referer.
//
// The referer is valid until returning from RequestHandler.
func (ctx *Context) Referer() string {
    return ctx.Request.Header.Get("referer")
}

// UserAgent returns User-Agent header value from the request.
func (ctx *Context) UserAgent() string {
    return ctx.Request.UserAgent()
}

// Serialize converts the given data into a byte array.
// If the data is neither a byte array nor a string, it will call fmt.Sprint to convert it into a string.
func Serialize(data interface{}) (bytes []byte, err error) {
    switch data.(type) {
    case []byte:
        return data.([]byte), nil
    case string:
        return []byte(data.(string)), nil
    default:
        if data != nil {
            return []byte(fmt.Sprint(data)), nil
        }
    }
    return nil, nil
}
