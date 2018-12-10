package main

import (
    "flag"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/zebra"
    //"github.com/tietang/logrus-prefixed-formatter"
    _ "net/http/pprof"
    "runtime"
    "strings"
)

type Bootstrap struct {
    bootstrap     string
    file          string
    zkUrls        string
    zkRoot        string
    consulAddress string
    consulRoot    string
    //
    run     bool
    start   bool
    restart bool
    stop    bool
    install bool
    status  bool
}

//var b = &Bootstrap{}
//--bootstrap|-b file|zk|consul
//--urls|-u 如果为zk|consul，指定连接字符串
//--path|-p 如果为zk|consul，指定配置的根路径
func Init(b *Bootstrap) {
    flag.StringVar(&b.bootstrap, "bootstrap", "file", "bootstrap类型，可选值为file|zk|consul")
    flag.StringVar(&b.file, "file", "proxy.ini", "指定配置文件，默认在当前目录下查找，如果为绝对路径则直接读取。")
    flag.StringVar(&b.zkUrls, "zkUrls", "127.0.0.1:2181", "多个地址用逗号分隔")
    flag.StringVar(&b.zkRoot, "zkRoot", "/zebra/bootstrap", "zNode root路径")
    flag.StringVar(&b.consulAddress, "address", "127.0.0.1:8500", "Consul Agent 地址")
    flag.StringVar(&b.consulRoot, "root", "zebra/bootstrap", "Consul key root路径")
    //flag.BoolVar(&b.run, "run", true, "控制台运行")
    //flag.BoolVar(&b.start, "start", false, "启动并后台运行服务")
    //flag.BoolVar(&b.restart, "restart", false, "重启")
    //flag.BoolVar(&b.stop, "stop", false, "停止")
    //flag.BoolVar(&b.install, "install", false, "从远程安装")
    //flag.BoolVar(&b.status, "status", false, "当前状态")

    //if b.isDaemon {
    //    cmd := exec.Command(os.Args[0], os.Args[1:]...)
    //    cmd.Start()
    //    log.Info(fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid))
    //    b.isDaemon = false
    //    os.Exit(0)
    //}
}
func logDump() {

}
func main() {
    logDump()
    b := &Bootstrap{}
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
    var server *zebra.ProxyServer
    if b.bootstrap == "file" {

        server = zebra.NewByFile(b.file)
    }
    if b.bootstrap == "zk" {
        urls := strings.Split(b.zkUrls, ",") //[]string{"127.0.0.1:2181"}
        root := b.zkRoot
        server = zebra.NewByZookeeper(root, urls)
    }
    if b.bootstrap == "consul" {
        address := b.consulAddress // "127.0.0.1:8500"
        root := b.consulRoot       //"zebra/bootstrap"
        server = zebra.NewByConsulKeyValue(address, root)
    }
    server.Start()

}
