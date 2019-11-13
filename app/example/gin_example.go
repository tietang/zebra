package main

//
//import (
//	"github.com/gin-gonic/gin"
//	"github.com/tietang/props/ini"
//	"github.com/tietang/zebra/router"
//	"strconv"
//)
//
//func main() {
//
//	file := "/Users/tietang/my/gitcode/r_app/src/github.com/tietang/zebra/app/proxy.ini"
//	conf := ini.NewIniFileConfigSource(file)
//	handler := router.NewUniversalHandler(conf)
//	r := gin.Default()
//	r.GET("/ping", func(c *gin.Context) {
//		c.JSON(200, gin.H{
//			"message": "pong",
//		})
//
//	})
//	r.Use(func(ctx *gin.Context) {
//		ok := handler.Handle(ctx.Writer, ctx.Request)
//		if !ok {
//			ctx.Next()
//		}
//	})
//	port := conf.GetIntDefault("server.port", 7980)
//	r.Run(":" + strconv.Itoa(port)) // listen and serve on 0.0.0.0:19001
//}
