package main

import (
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/router"
    "net/http"
)

func main() {
    file := "/Users/tietang/Documents/git.oschina/r_app/src/github.com/tietang/zebra/example/proxy.ini"
    conf := kvs.NewIniFileConfigSource(file)
    handler := router.NewUniversalHandler(conf)

    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ok := handler.Handle(writer, request)
        if !ok {
            //
        }
    })

    http.ListenAndServe(":8080", nil)
}
