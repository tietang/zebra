package main

import (
    "github.com/go-martini/martini"
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/router"
    "net/http"
    "strconv"
)

func main() {
    file := "/Users/tietang/my/gitcode/r_app/src/github.com/tietang/zebra/app/proxy.ini"
    conf := kvs.NewIniFileConfigSource(file)
    handler := router.NewUniversalHandler(conf)
    m := martini.Classic()
    m.Get("/", func() string {
        return "Hello world!"
    })
    m.Use(func(c martini.Context, res http.ResponseWriter, req *http.Request) {

        ok := handler.Handle(res, req)
        if !ok {
            c.Next()
        }
    })
    port := conf.GetIntDefault("server.port", 7980)
    m.RunOnAddr(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:19001
}
