package main

import (
	"github.com/kataras/iris/v12"
	"github.com/tietang/props/ini"
	"github.com/tietang/zebra/router"
	"strconv"

	zutils "github.com/tietang/zebra/utils"
)

func main() {
	file := "/Users/tietang/my/gitcode/r_app/src/github.com/tietang/zebra/app/proxy.ini"
	conf := ini.NewIniFileConfigSource(file)
	handler := router.NewUniversalHandler(conf)
	app := iris.New()
	// Resource:  http://localhost:8080
	app.Any("/{asset:path}", func(ctx iris.Context) {
		ok := handler.Handle(ctx.ResponseWriter(), ctx.Request())
		if !ok {
			ctx.Next()
		}
	})
	// Method:   GET
	// Resource: http://localhost:8080/
	app.Handle("GET", "/", func(ctx iris.Context) {
		ctx.HTML("<b>Hello world!</b>")
	})
	app.Handle("GET", "/favicon.ico", func(ctx iris.Context) {
		ctx.Binary(zutils.ICON_DATA)
	})
	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://context:8080/ping
	app.Get("/ping", func(ctx iris.Context) {
		ctx.WriteString("pong")
	})

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx iris.Context) {
		ctx.JSON(iris.Map{"message": "Hello iris web framework."})
	})

	port := conf.GetIntDefault("server.port", 7980)
	app.Run(iris.Addr(":" + strconv.Itoa(port))) // listen and serve on 0.0.0.0:8080
}
