package main

import (
	"github.com/kataras/iris"
	"github.com/kataras/iris/context"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/router"
	"strconv"

	"github.com/mattn/go-colorable"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	zutils "github.com/tietang/zebra/utils"
)

func init() {
	//formatter := &log.TextFormatter{}
	//formatter.ForceColors = true
	//formatter.DisableColors = false
	//formatter.FullTimestamp = true
	//formatter.TimestampFormat = "2006-01-02.15:04:05.999999"
	//

	//formatter := &prefixed.TextFormatter{}
	formatter := &utils.TextFormatter{}
	formatter.ForceColors = true
	formatter.DisableColors = false
	formatter.FullTimestamp = true
	formatter.ForceFormatting = true
	//formatter.EnableFuncNameLog = false
	formatter.SetColorScheme(&utils.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "red",
		PanicLevelStyle: "red",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan+b",
		TimestampStyle:  "black+h",
	})
	formatter.TimestampFormat = "2006-01-02.15:04:05.999999"

	log.SetFormatter(formatter)
	log.SetOutput(colorable.NewColorableStdout())
	//log.SetOutput(os.Stdout) propfile
	log.SetLevel(log.DebugLevel)

}

func main() {
	file := "/Users/tietang/my/gitcode/r_app/src/github.com/tietang/zebra/app/proxy.ini"
	conf := kvs.NewIniFileConfigSource(file)
	handler := router.NewUniversalHandler(conf)
	app := iris.New()
	// Resource:  http://localhost:8080
	app.Any("/{asset:path}", func(ctx context.Context) {
		ok := handler.Handle(ctx.ResponseWriter(), ctx.Request())
		if !ok {
			ctx.Next()
		}
	})
	// Method:   GET
	// Resource: http://localhost:8080/
	app.Handle("GET", "/", func(ctx context.Context) {
		ctx.HTML("<b>Hello world!</b>")
	})
	app.Handle("GET", "/favicon.ico", func(ctx context.Context) {
		ctx.Binary(zutils.ICON_DATA)
	})
	// same as app.Handle("GET", "/ping", [...])
	// Method:   GET
	// Resource: http://context:8080/ping
	app.Get("/ping", func(ctx context.Context) {
		ctx.WriteString("pong")
	})

	// Method:   GET
	// Resource: http://localhost:8080/hello
	app.Get("/hello", func(ctx context.Context) {
		ctx.JSON(context.Map{"message": "Hello iris web framework."})
	})

	port := conf.GetIntDefault("server.port", 7980)
	app.Run(iris.Addr(":" + strconv.Itoa(port))) // listen and serve on 0.0.0.0:8080
}
