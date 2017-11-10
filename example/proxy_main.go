package main

import "github.com/tietang/zebra"

//type Handler struct {
//	Count int
//}

//func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
//	w.Write([]byte(`
////	Go port of Coda Hale's Metrics library: https://github.com/dropwizard/metrics.
//	`))
//	//	time.Sleep(time.Millisecond * 100)
//	url := req.URL
//	uri := req.RequestURI
//	w.Write([]byte(strconv.Itoa(h.Count) + " " + uri + "  " + url.Path + "  |  " + url.RawPath + "  |  " + url.RawQuery))
//	h.Count++
//}

func main() {
    eurekaUrl := "http://127.0.0.1:8761/eureka"

    //	http.HandleFunc("/info", func(w http.ResponseWriter, req *http.Request) {
    //		w.Write([]byte("{}"))
    //	})
    //	http.ListenAndServe(":8080", &Handler{Count: 0})
    //	http.ListenAndServe(":8080", nil)
    //r.AddRoute(Route{Source: "/app1/v1/user", App: "app1", Target: "/v1/user", StripPrefix: false})

    routes := []string{"/app1/v1/user,app1,/v1/user", "/app1/info,app1,/info"}
    proxy := zuul.NewHttpProxyServer(eurekaUrl, routes)
    proxy.DefaultRun()
}
