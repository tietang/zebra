package main

import (
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/router"
    "golang.org/x/net/http2"
    "net/http"
)

func main() {
    file := "/Users/tietang/Documents/git.oschina/r_app/src/github.com/tietang/zebra/example/proxy.ini"
    conf := kvs.NewIniFileConfigSource(file)
    handler := router.NewUniversalHandler(conf)
    var server http.Server

    http2.VerboseLogs = true
    server.Addr = ":8080"

    http2.ConfigureServer(&server, &http2.Server{})

    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        handler.Handle(writer, request)
    })
    server.ListenAndServe()
}
