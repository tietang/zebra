package main

import (
    "github.com/gin-gonic/gin"
    "github.com/mattn/go-colorable"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/go-utils"
    "github.com/tietang/props/kvs"
    "github.com/tietang/zebra/router"
    "strconv"
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
    //formatter.EnableLogFuncName = false
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
    r := gin.Default()
    r.GET("/ping", func(c *gin.Context) {
        c.JSON(200, gin.H{
            "message": "pong",
        })

    })
    r.Use(func(ctx *gin.Context) {
        ok := handler.Handle(ctx.Writer, ctx.Request)
        if !ok {
            ctx.Next()
        }
    })
    port := conf.GetIntDefault("server.port", 7980)
    r.Run(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:19001
}
