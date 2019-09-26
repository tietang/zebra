package zebra

import (
	"errors"
	"fmt"
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
	"io"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path"
	"path/filepath"
	"runtime"
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
var lfh *utils.LineNumLogrusHook

func init() {
	// 定义日志格式
	formatter = &prefixed.TextFormatter{}
	//设置高亮显示的色彩样式
	formatter.ForceColors = true
	formatter.DisableColors = false
	formatter.ForceFormatting = true
	formatter.SetColorScheme(&prefixed.ColorScheme{
		InfoLevelStyle:  "green",
		WarnLevelStyle:  "yellow",
		ErrorLevelStyle: "red",
		FatalLevelStyle: "41",
		PanicLevelStyle: "41",
		DebugLevelStyle: "blue",
		PrefixStyle:     "cyan",
		TimestampStyle:  "37",
	})
	//开启完整时间戳输出和时间戳格式
	formatter.FullTimestamp = true
	//设置时间格式
	formatter.TimestampFormat = "2006-01-02.15:04:05.000000"
	//设置日志formatter
	log.SetFormatter(formatter)
	log.SetOutput(colorable.NewColorableStdout())
	//日志级别，通过环境变量来设置
	// 后期可以变更到配置中来设置

	if os.Getenv("log.debug") == "true" {
		log.SetLevel(log.DebugLevel)
	}
	//开启调用函数、文件、代码行信息的输出
	log.SetReportCaller(true)
	//设置函数、文件、代码行信息的输出的hook
	SetLineNumLogrusHook()

}

func SetLineNumLogrusHook() {
	lfh = utils.NewLineNumLogrusHook()
	lfh.EnableFileNameLog = true
	lfh.EnableFuncNameLog = true
	log.AddHook(lfh)
}

//将滚动日志writer共享给iris glog output
var log_writer io.Writer

//初始化log配置，配置logrus日志文件滚动生成和
func InitLog(conf kvs.ConfigSource) {
	//设置日志输出级别
	level, err := log.ParseLevel(conf.GetDefault("log.level", "info"))
	if err != nil {
		level = log.InfoLevel
	}
	log.SetLevel(level)
	if conf.GetBoolDefault("log.enableLineLog", true) {
		lfh.EnableFileNameLog = true
		lfh.EnableFuncNameLog = true
	} else {
		lfh.EnableFileNameLog = false
		lfh.EnableFuncNameLog = false
	}

	//配置日志输出目录
	logDir := conf.GetDefault("log.dir", "./logs")
	logTestDir, err := conf.Get("log.test.dir")
	if err == nil {
		logDir = logTestDir
	}
	logPath := logDir //+ "/logs"
	logFilePath, _ := filepath.Abs(logPath)
	log.Infof("log dir: %s", logFilePath)
	logFileName := conf.GetDefault("log.file.name", "red-envelop")
	maxAge := conf.GetDurationDefault("log.max.age", time.Hour*24)
	rotationTime := conf.GetDurationDefault("log.rotation.time", time.Hour*1)
	os.MkdirAll(logPath, os.ModePerm)

	baseLogPath := path.Join(logPath, logFileName)
	//设置滚动日志输出writer
	writer, err := rotatelogs.New(
		strings.TrimSuffix(baseLogPath, ".log")+".%Y%m%d%H.log",
		rotatelogs.WithLinkName(baseLogPath),      // 生成软链，指向最新日志文件
		rotatelogs.WithMaxAge(maxAge),             // 文件最大保存时间
		rotatelogs.WithRotationTime(rotationTime), // 日志切割时间间隔
	)
	if err != nil {
		log.Errorf("config local file system logger error. %+v", err)
	}
	//设置日志文件输出的日志格式
	formatter := &log.TextFormatter{}
	formatter.CallerPrettyfier = func(frame *runtime.Frame) (function string, file string) {
		function = frame.Function
		dir, filename := path.Split(frame.File)
		f := path.Base(dir)
		return function, fmt.Sprintf("%s/%s:%d", f, filename, frame.Line)
	}
	lfHook := lfshook.NewHook(lfshook.WriterMap{
		log.DebugLevel: writer, // 为不同级别设置不同的输出目的
		log.InfoLevel:  writer,
		log.WarnLevel:  writer,
		log.ErrorLevel: writer,
		log.FatalLevel: writer,
		log.PanicLevel: writer,
	}, formatter)

	log.AddHook(lfHook)
	//
	log_writer = writer

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

//func (p *ProxyServer) RegisterSource(t interface{}) {
//
//}

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
