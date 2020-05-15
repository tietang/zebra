package proxy

import (
	"encoding/json"
	"fmt"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/infra"
)

type AdminEndpoints struct {
	server *HttpProxyServer
	conf   kvs.ConfigSource
}

func NewAdminEndpoints(server *HttpProxyServer, conf kvs.ConfigSource) *AdminEndpoints {
	return &AdminEndpoints{
		server: server,
		conf:   conf,
	}
}

//http://localhost:19091/admin/gray
//app=RESK&version=1.1&type=addrs&key=&values=10.180.7.112&ranges=&enabled=true
func (a *AdminEndpoints) AdminGrayConfEndpoint(ctx *infra.Context) error {
	//if faviconIconData == nil || len(faviconIconData) == 0 {
	//    Path := h.conf.GetDefault(KEY_FAVICON_ICO_PATH, "favicon.ico")
	//    data, err := ioutil.ReadFile(Path)
	//    if err != nil {
	//        return err
	//    }
	//    faviconIconData = data
	//}

	//q := ctx.URL().Query()
	err := ctx.Request.ParseForm()
	fmt.Println(err)
	q := ctx.Request.PostForm
	fmt.Println(q)
	app := q.Get("app")
	version := q.Get("version")
	enabled := q.Get("enabled")
	typ := q.Get("type")
	props := make(map[string]string)
	if enabled == "true" {
		props[fmt.Sprintf("traffic.cond.%s.enabled", app)] = enabled
		props[fmt.Sprintf("traffic.cond.%s.version", app)] = version
	}
	if typ == "addrs" {
		addrs := q.Get("values")
		props[fmt.Sprintf("traffic.cond.%s.client.addrs", app)] = addrs
	}
	if typ == "header" {
		key := q.Get("key")
		values := q.Get("values")
		ranges := q.Get("ranges")
		props[fmt.Sprintf("traffic.cond.%s.header.key", app)] = key
		props[fmt.Sprintf("traffic.cond.%s.header.values", app)] = values
		props[fmt.Sprintf("traffic.cond.%s.header.ranges", app)] = ranges

	}
	if typ == "cookie" {
		key := q.Get("key")
		values := q.Get("values")
		ranges := q.Get("ranges")
		props[fmt.Sprintf("traffic.cond.%s.cookie.key", app)] = key
		props[fmt.Sprintf("traffic.cond.%s.cookie.values", app)] = values
		props[fmt.Sprintf("traffic.cond.%s.cookie.ranges", app)] = ranges

	}
	if typ == "urlParams" {
		key := q.Get("key")
		values := q.Get("values")
		ranges := q.Get("ranges")
		props[fmt.Sprintf("traffic.cond.%s.urlParams.key", app)] = key
		props[fmt.Sprintf("traffic.cond.%s.urlParams.values", app)] = values
		props[fmt.Sprintf("traffic.cond.%s.urlParams.ranges", app)] = ranges

	}
	if typ == "form" {
		key := q.Get("key")
		values := q.Get("values")
		ranges := q.Get("ranges")
		props[fmt.Sprintf("traffic.cond.%s.form.key", app)] = key
		props[fmt.Sprintf("traffic.cond.%s.form.values", app)] = values
		props[fmt.Sprintf("traffic.cond.%s.form.ranges", app)] = ranges
	}
	fmt.Println("props:", props)
	a.conf.SetAll(props)
	//
	d, _ := json.Marshal(props)

	ctx.WriteData(d)

	return nil
}
