package proxy

import (
	"github.com/tietang/zebra"
	"github.com/tietang/zebra/infra"
	"strings"
)

var bootstrap *zebra.Bootstrap

func SetBootstrap(b *zebra.Bootstrap) {
	bootstrap = b
}

type ProxyServerStarter struct {
	infra.BaseStarter
	proxy *ProxyServer
}

func (i *ProxyServerStarter) Init(ctx infra.StarterContext) {
	if bootstrap.Bootstrap == "file" {
		i.proxy = NewByFile(bootstrap.File)
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

}

func (i *ProxyServerStarter) Start(ctx infra.StarterContext) {
	i.proxy.Start()
}
func (i *ProxyServerStarter) StartBlocking() bool {
	return true
}
