package zebra

import (
    "errors"
    "github.com/lestrrat/go-file-rotatelogs"
    "github.com/mattn/go-colorable"
    "github.com/rifflock/lfshook"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/go-utils"
    "github.com/tietang/props/consul"
    "github.com/tietang/props/ini"
    "github.com/tietang/props/kvs"
    "github.com/tietang/props/zk"
    "github.com/tietang/zebra/router"
    "github.com/x-cray/logrus-prefixed-formatter"
    "net/http"
    _ "net/http/pprof"
    "os"
    "path"
    "path/filepath"
    "strings"
    "time"
)

//bootstrap
const (
    //key
    KEY_BOOTSTRAP_ZK_ENABLED     = "bootstrap.zk.enabled"
    KEY_BOOTSTRAP_ZK_URLS        = "bootstrap.zk.urls"
    KEY_BOOTSTRAP_ZK_TIMEOUT     = "bootstrap.zk.timeout"
    KEY_BOOTSTRAP_ZK_ROOT        = "bootstrap.zk.root"
    KEY_BOOTSTRAP_CONSUL_ENABLED = "bootstrap.consul.enabled"
    KEY_BOOTSTRAP_CONSUL_ADDRESS = "bootstrap.consul.address"
    KEY_BOOTSTRAP_CONSUL_ROOT    = "bootstrap.consul.root"
    //default value
    DEFAULT_BOOTSTRAP_ZK_ENABLED     = "false"
    DEFAULT_BOOTSTRAP_ZK_URLS        = "127.0.0.1:2181"
    DEFAULT_BOOTSTRAP_ZK_TIMEOUT     = time.Duration(10000) * time.Millisecond
    DEFAULT_BOOTSTRAP_ZK_ROOT        = "/zebra/bootstrap"
    DEFAULT_BOOTSTRAP_CONSUL_ENABLED = "false"
    DEFAULT_BOOTSTRAP_CONSUL_ADDRESS = "127.0.0.1:8500"
    DEFAULT_BOOTSTRAP_CONSUL_ROOT    = "zebra/bootstrap"
)

var formatter *prefixed.TextFormatter

func init() {
    //formatter := &log.TextFormatter{}
    //formatter.ForceColors = true
    //formatter.DisableColors = false
    //formatter.FullTimestamp = true
    //formatter.TimestampFormat = "2006-01-02.15:04:05.999999"
    //

    //formatter := &prefixed.TextFormatter{}
    formatter = &prefixed.TextFormatter{}
    formatter.ForceColors = true
    formatter.DisableColors = false
    formatter.FullTimestamp = true
    formatter.ForceFormatting = true
    //formatter.EnableLogFuncName = false
    formatter.SetColorScheme(&prefixed.ColorScheme{
        InfoLevelStyle:  "green",
        WarnLevelStyle:  "yellow",
        ErrorLevelStyle: "red",
        FatalLevelStyle: "red",
        PanicLevelStyle: "red",
        DebugLevelStyle: "blue",
        PrefixStyle:     "cyan+b",
        TimestampStyle:  "black+h",
    })
    formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
    log.SetFormatter(formatter)
    log.SetOutput(colorable.NewColorableStdout())
    //log.SetOutput(os.Stdout) propfile
    if os.Getenv("log.debug") == "true" {
        log.SetLevel(log.DebugLevel)
    }
    //
    logPath := "./logs/"
    logFileName := "zebra"
    maxAge := time.Hour * 24
    rotationTime := time.Hour
    os.MkdirAll(logPath, os.ModePerm)

    baseLogPath := path.Join(logPath, logFileName)
    writer, err := rotatelogs.New(
        baseLogPath+".%Y%m%d%H%M.log",
        rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
        rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
        rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
    )
    if err != nil {
        log.Errorf("config local file system logger error. %+v", err)
    }
    lfHook := lfshook.NewHook(lfshook.WriterMap{
        log.DebugLevel: writer, // 为不同级别设置不同的输出目的
        log.InfoLevel:  writer,
        log.WarnLevel:  writer,
        log.ErrorLevel: writer,
        log.FatalLevel: writer,
        log.PanicLevel: writer,
    })
    log.AddHook(lfHook)
    log.AddHook(utils.NewLogLineNumHook())
}

type ProxyServer struct {
    ConfigFilePath  string
    conf            kvs.ConfigSource
    HttpProxyServer *HttpProxyServer
}

func NewByFile(fileName string) *ProxyServer {
    //s := time.Now().UnixNano()
    file := fileName
    if !filepath.IsAbs(fileName) {
        dir, err := os.Getwd()
        utils.Panic(err)
        file = filepath.Join(dir, fileName)
    }

    isExists, err := utils.PathExists(file)
    utils.Panic(err)
    if !isExists {
        //panic(errors.New("file is not exists: " + file))
        log.Error(errors.New("file is not exists: "+file), err)
    }

    conf := kvs.NewEmptyCompositeConfigSource()
    //优先级的顺序为consul> zookeeper > properties
    ext := filepath.Ext(file)
    var rootConf kvs.ConfigSource
    if strings.Contains(ext, "prop") {
        rootConf = kvs.NewPropertiesConfigSource(file)
    } else {

        rootConf = ini.NewIniFileConfigSource(file)

    }

    //

    consulAddress := rootConf.GetDefault(KEY_BOOTSTRAP_CONSUL_ADDRESS, DEFAULT_BOOTSTRAP_CONSUL_ADDRESS)
    if rootConf.GetBoolDefault(KEY_BOOTSTRAP_CONSUL_ENABLED, false) && err == nil && consulAddress != "" {
        root := rootConf.GetDefault(KEY_BOOTSTRAP_CONSUL_ROOT, DEFAULT_BOOTSTRAP_CONSUL_ROOT)
        consulConf := consul.NewConsulKeyValueConfigSource(consulAddress, root)
        conf.Add(consulConf)
    }
    //
    zkUrls := rootConf.GetDefault(KEY_BOOTSTRAP_ZK_URLS, DEFAULT_BOOTSTRAP_ZK_URLS)
    if rootConf.GetBoolDefault(KEY_BOOTSTRAP_ZK_ENABLED, false) && err == nil && zkUrls != "" {
        contexts := []string{rootConf.GetDefault(KEY_BOOTSTRAP_ZK_ROOT, DEFAULT_BOOTSTRAP_ZK_ROOT)}
        connStr := strings.Split(zkUrls, ",")
        timeout := rootConf.GetDurationDefault(KEY_BOOTSTRAP_ZK_TIMEOUT, DEFAULT_BOOTSTRAP_ZK_TIMEOUT)
        zkConf := zk.NewZookeeperCompositeConfigSource(contexts, connStr, timeout)
        conf.Add(zkConf)
    }

    conf.Add(rootConf)
    p := New(conf)
    //t := time.Now().UnixNano() - s
    //fmt.Println(t, "ns ", strconv.Itoa(int(t/int64(time.Millisecond))), "ms")
    return p
}

func NewByZookeeper(root string, urls []string) *ProxyServer {
    conf := zk.NewZookeeperCompositeConfigSource([]string{root}, urls, time.Second*6)
    p := New(conf)
    return p
}

func NewByConsulKeyValue(address, root string) *ProxyServer {
    conf := consul.NewConsulKeyValueConfigSource(address, root)
    p := New(conf)
    return p
}

func New(conf kvs.ConfigSource) *ProxyServer {

    if conf.GetBoolDefault("server.log.debug", false) {
        log.SetLevel(log.DebugLevel)
    }

    if conf.GetBoolDefault("server.log.color", false) {
        formatter.ForceColors = false
        formatter.DisableColors = true
    }

    server := NewHttpProxyServer(conf)
    p := &ProxyServer{
        //ConfigFilePath:  file,
        conf:            conf,
        HttpProxyServer: server,
    }

    return p
}

func (p *ProxyServer) Start() {
    //p.Use(func(ctx *proxy.Context) error {
    //    ctx.Write([]byte("hello world"))
    //    ctx.Next()
    //    return nil
    //})
    port := p.conf.GetDefault("server.admin.port", "60001")
    go func() {
        log.Info("listened and served admin port: ", port)
        log.Info(http.ListenAndServe(":"+port, nil))
    }()
    log.Info("init router...")
    p.setupRouter()
    log.Info("starting...")
    p.run()
}

func (p *ProxyServer) RegisterSource(t interface{}) {

}

func (p *ProxyServer) run() {
    p.HttpProxyServer.DefaultRun()
}

// Use appends Handler(s) to the current Party's routes and child routes.
// If the current Party is the root, then it registers the middleware to all child Parties' routes too.
func (r *ProxyServer) Use(handlers ...Handler) {
    r.HttpProxyServer.Use(handlers...)
}

func (h *ProxyServer) setupRouter() {

    handler := router.NewUniversalHandler(h.conf)

    h.HttpProxyServer.health.AddAll(handler.Router.GetHealthCheckers())

    h.Use(func(context *Context) error {
        handler.Handle(context.ResponseWriter, context.Request)
        return nil
    })

}
