package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/consul"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/zk"
	"github.com/tietang/zebra"
	"github.com/tietang/zebra/infra"
	"github.com/tietang/zebra/infra/base"
	"github.com/tietang/zebra/proxy"
	"path/filepath"
	"time"

	//"github.com/tietang/logrus-prefixed-formatter"
	_ "net/http/pprof"
	"runtime"
	"strings"
)

//var b = &Bootstrap{}
//--bootstrap|-b file|zk|consul
//--urls|-u 如果为zk|consul，指定连接字符串
//--path|-p 如果为zk|consul，指定配置的根路径
func Init(b *zebra.Bootstrap) {
	flag.StringVar(&b.Bootstrap, "bootstrap", "file", "bootstrap类型，可选值为file|zk|consul")
	flag.StringVar(&b.File, "file", "proxy.ini", "指定配置文件，默认在当前目录下查找，如果为绝对路径则直接读取。")
	flag.StringVar(&b.ZkUrls, "zkUrls", "127.0.0.1:2181", "多个地址用逗号分隔")
	flag.StringVar(&b.ZkRoot, "zkRoot", "/zebra/bootstrap", "zNode root路径")
	flag.StringVar(&b.ConsulAddress, "address", "127.0.0.1:8500", "Consul Agent 地址")
	flag.StringVar(&b.ConsulRoot, "root", "zebra/bootstrap", "Consul key root路径")
}

func main() {
	logDump()
	b := &zebra.Bootstrap{}
	Init(b)
	if !flag.Parsed() {
		flag.Parse()
	}

	runtime.GOMAXPROCS(runtime.NumCPU())
	//for i, param := range flag.Args() {
	//    fmt.Printf("#%d    :%s\n", i, param)
	//}
	log.WithField("bootstrap", b).Debug()
	//fmt.Println(b)

	//f, err := os.Create("trace.out")
	//if err != nil {
	//    panic(err)
	//}
	//defer f.Close()
	//
	//err = trace.Start(f)
	//if err != nil {
	//    panic(err)
	//}
	//defer trace.Stop()
	app := infra.New(conf)
	app.Start()

	start(b)
}

//
//
//func start(b *Bootstrap) {
//	var server *proxy.ProxyServer
//	if b.bootstrap == "file" {
//
//		server = proxy.NewByFile(b.file)
//	}
//	if b.bootstrap == "zk" {
//		urls := strings.Split(b.zkUrls, ",") //[]string{"127.0.0.1:2181"}
//		root := b.zkRoot
//		server = proxy.NewByZookeeper(root, urls)
//	}
//	if b.bootstrap == "consul" {
//		address := b.consulAddress // "127.0.0.1:8500"
//		root := b.consulRoot       //"zebra/bootstrap"
//		server = proxy.NewByConsulKeyValue(address, root)
//	}
//	server.Start()
//
//}
