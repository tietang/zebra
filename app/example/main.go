package main

import (
	"flag"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/zebra/proxy"

	//"github.com/tietang/logrus-prefixed-formatter"
	_ "net/http/pprof"
	"runtime"
	"strings"
	//_  "github.com/codyguo/godaemon"
	//"os"
	//"os/exec"
	"os"
	"syscall"
	//仅导入，包的init方法被自动调用，嵌入daemon功能
	_ "github.com/tietang/godaemon"
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

var b = &Bootstrap{}

//--bootstrap|-b file|zk|consul
//--urls|-u 如果为zk|consul，指定连接字符串
//--path|-p 如果为zk|consul，指定配置的根路径
func init() {
	flag.StringVar(&b.bootstrap, "bootstrap", "file", "bootstrap类型，可选值为file|zk|consul")
	flag.StringVar(&b.file, "file", "proxy.ini", "指定配置文件，默认在当前目录下查找，如果为绝对路径则直接读取。")
	flag.StringVar(&b.zkUrls, "zkUrls", "127.0.0.1:2181", "多个地址用逗号分隔")
	flag.StringVar(&b.zkRoot, "zkRoot", "/zebra/bootstrap", "zNode root路径")
	flag.StringVar(&b.consulAddress, "address", "127.0.0.1:8500", "Consul Agent 地址")
	flag.StringVar(&b.consulRoot, "root", "zebra/bootstrap", "Consul key root路径")
	flag.BoolVar(&b.run, "run", true, "控制台运行")
	flag.BoolVar(&b.start, "start", false, "启动并后台运行服务")
	flag.BoolVar(&b.restart, "restart", false, "重启")
	flag.BoolVar(&b.stop, "stop", false, "停止")
	flag.BoolVar(&b.install, "install", false, "从远程安装")
	flag.BoolVar(&b.status, "status", false, "当前状态")
	if !flag.Parsed() {
		flag.Parse()
	}
	//if b.isDaemon {
	//    cmd := exec.Command(os.Args[0], os.Args[1:]...)
	//    cmd.Start()
	//    log.Info(fmt.Printf("%s [PID] %d running...\n", os.Args[0], cmd.Process.Pid))
	//    b.isDaemon = false
	//    os.Exit(0)
	//}
}

func main() {
	logFile, err := os.OpenFile("./dump.log", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0660)
	if err != nil {
		log.Println("服务启动出错", "打开异常日志文件失败", err)
		return
	}
	// 将进程标准出错重定向至文件，进程崩溃时运行时将向该文件记录协程调用栈信息,
	// linux系统中的dup系统调用在windows系统及mac系统下暂时不可用
	syscall.Dup2(int(logFile.Fd()), int(os.Stderr.Fd()))

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
	var server *proxy.ProxyServer
	if b.bootstrap == "file" {

		server = proxy.NewByFile(b.file)
	}
	if b.bootstrap == "zk" {
		urls := strings.Split(b.zkUrls, ",") //[]string{"127.0.0.1:2181"}
		root := b.zkRoot
		server = proxy.NewByZookeeper(root, urls)
	}
	if b.bootstrap == "consul" {
		address := b.consulAddress // "127.0.0.1:8500"
		root := b.consulRoot       //"zebra/bootstrap"
		server = proxy.NewByConsulKeyValue(address, root)
	}
	server.Start()

}
