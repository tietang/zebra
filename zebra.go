package zebra

import (
	"gitee.com/tietang/terrace-go/base"
	"gitee.com/tietang/terrace-go/boot"
	"github.com/tietang/zebra/proxy"
)

func init() {
	boot.Register(&base.PropsStarter{})
	boot.Register(&boot.WebApiStarter{})
	boot.Register(&boot.HookStarter{})
	boot.Register(&proxy.ProxyServerStarter{})
}
