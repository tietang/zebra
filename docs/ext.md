
#自定义或在主流web framework上扩展

##Installation

需要[go](<https://golang.org/dl/>)环境,版本1.8+。

安装zebra：

```
go get -u github.com/tietang/zebra/router
```

安装依赖：

```
go get -u github.com/tietang/props
go get -u github.com/tietang/go-utils
```



##原生http例子

`props.ConfigSource`使用方法请参考[props](<https://github.com/tietang/props>)，支持内存properties文件或硬编码、zookeeper key/value、consul key/value 3中配置支持。

http11:


```golang
 var conf kvs.ConfigSource
 handler := router.NewUniversalHandler(conf)

 http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
    ok := handler.Handle(writer, request)
        if !ok {
            //
        }
 })

 http.ListenAndServe(":8080", nil)
```
http2:

```golang

    handler := router.NewUniversalHandler(conf)
    var server http.Server

    http2.VerboseLogs = true
    server.Addr = ":8080"

    http2.ConfigureServer(&server, &http2.Server{})

    http.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
        ok := handler.Handle(writer, request)
        if !ok {
            //
        }
    })
    server.ListenAndServe() 
```

主要扩展点是调用`UniversalHandler.Handle(w http.ResponseWriter, req *http.Request)`来扩展。


注意：要通过`router.NewUniversalHandler(conf)`提供的New函数来构造`UniversalHandler`，不建议通过`&UniversalHandler{}`方式自己来new。

```golang
//创建一个配置，并读取配置，可以参考https://github.com/tietang/props
var conf kvs.ConfigSource
handler := router.NewUniversalHandler(conf)
func (h *UniversalHandler) Handle(w http.ResponseWriter, req *http.Request) bool
```

## gin扩展例子

```golang
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
	
```

## iris

```golang
handler := router.NewUniversalHandler(conf)
	app := iris.New()
	// Resource:  http://localhost:8080
	app.Any("/{asset:path}", func(ctx iris.Context) {
		ok := handler.Handle(ctx.ResponseWriter(), ctx.Request())
		if !ok {
			ctx.Next()
		}
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
		ctx.JSON(context.Map{"message": "Hello iris web framework."})
	})
	

	port := conf.GetIntDefault("server.port", 7980)
	app.Run(iris.Addr(":" + strconv.Itoa(port))) // listen and serve on 0.0.0.0:8080
```
## martini


```golang
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
```