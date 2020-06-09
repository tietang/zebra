package proxy

import (
	"flag"
	"gitee.com/tietang/terrace-go/boot"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"runtime"
	"strings"
)

type Bootstrap struct {
	Bootstrap     string
	File          string
	ZkUrls        string
	ZkRoot        string
	ConsulAddress string
	ConsulRoot    string
}

var bootstrap *Bootstrap

func SetBootstrap(b *Bootstrap) {
	bootstrap = b
}

type ProxyServerStarter struct {
	boot.BaseStarter
	proxy *ProxyServer
}

func (i *ProxyServerStarter) Init(ctx boot.StarterContext) {

	logDump()
	bootstrap = &Bootstrap{}
	InitBootstrapByArgs(bootstrap)
	if !flag.Parsed() {
		flag.Parse()
	}
	runtime.GOMAXPROCS(runtime.NumCPU())
	//for i, param := range flag.Args() {
	//    fmt.Printf("#%d    :%s\n", i, param)
	//}
	log.WithField("bootstrap", bootstrap).Debug()
	if bootstrap.Bootstrap == "file" {
		//i.proxy = NewByFile(bootstrap.File)
		i.proxy = New(ctx.Props())
	}
	if bootstrap.Bootstrap == "zk" {
		urls := strings.Split(bootstrap.ZkUrls, ",") //[]string{"127.0.0.1:2181"}
		root := bootstrap.ZkRoot
		i.proxy = NewByZookeeper(root, urls)
	}
	if bootstrap.Bootstrap == "consul" {
		address := bootstrap.ConsulAddress // "127.0.0.1:8500"
		root := bootstrap.ConsulRoot       //"zebra/bootstrap"
		i.proxy = NewByConsulKeyValue(address, root)
	}
	if i.proxy == nil {
		log.Panic("proxy server is nil.")
	}
	i.proxy.conf = ctx.Props()

}

func (i *ProxyServerStarter) Start(ctx boot.StarterContext) {
	i.proxy.Start()
}
func (i *ProxyServerStarter) StartBlocking() bool {
	return true
}

func InitBootstrapByFile(conf kvs.ConfigSource, b *Bootstrap) {
	//	##默认为properties配置，即为本文件的配置；如果配置了bootstrap.zk.enabled和bootstrap.consul.enabled值为true时，
	//	则zookeeper和consul配置也会起效，并且具有优先获取权，如果找不到再从properties中找。
	//	##优先级的顺序为consul> zookeeper > properties[file]
	if conf.GetBoolDefault("bootstrap.consul.enabled", false) {
		b.Bootstrap = "consul"
		b.ConsulAddress = conf.GetDefault("bootstrap.consul.address", "127.0.0.1:8500")
		b.ConsulRoot = conf.GetDefault("bootstrap.consul.root", "zebra/bootstrap")
	} else if conf.GetBoolDefault("bootstrap.zk.enabled", false) {
		b.Bootstrap = "zk"
		b.ZkUrls = conf.GetDefault("bootstrap.zk.urls", "127.0.0.1:2181")
		b.ZkRoot = conf.GetDefault("bootstrap.zk.root", "/zebra/bootstrap")
	} else {
		b.Bootstrap = "file"
	}

}

//var b = &Bootstrap{}
//--bootstrap|-b file|zk|consul
//--urls|-u 如果为zk|consul，指定连接字符串
//--path|-p 如果为zk|consul，指定配置的根路径
func InitBootstrapByArgs(b *Bootstrap) {
	flag.StringVar(&b.Bootstrap, "bootstrap", "file", "bootstrap类型，可选值为file|zk|consul")
	flag.StringVar(&b.File, "file", "proxy.ini", "指定配置文件，默认在当前目录下查找，如果为绝对路径则直接读取。")
	flag.StringVar(&b.ZkUrls, "zkUrls", "127.0.0.1:2181", "多个地址用逗号分隔")
	flag.StringVar(&b.ZkRoot, "zkRoot", "/zebra/bootstrap", "zNode root路径")
	flag.StringVar(&b.ConsulAddress, "address", "127.0.0.1:8500", "Consul Agent 地址")
	flag.StringVar(&b.ConsulRoot, "root", "zebra/bootstrap", "Consul key root路径")
}
