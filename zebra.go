package zebra

import (
	"github.com/tietang/zebra/infra"
	"github.com/tietang/zebra/infra/base"
	"github.com/tietang/zebra/proxy"
)

func init() {
	infra.Register(&base.PropsStarter{})
	infra.Register(&infra.WebApiStarter{})
	infra.Register(&base.HookStarter{})
	infra.Register(&proxy.ProxyServerStarter{})
}

type Bootstrap struct {
	Bootstrap     string
	File          string
	ZkUrls        string
	ZkRoot        string
	ConsulAddress string
	ConsulRoot    string
}
