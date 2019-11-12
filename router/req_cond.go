package router

import (
	"fmt"
	"github.com/tietang/props/kvs"
	"strconv"
	"strings"
)

var RequestConditions map[string]RequestCondition

func init() {
	RequestConditions = make(map[string]RequestCondition)
	register(new(ClientAddrRequestCondition))
	register(new(HeaderValueRequestCondition))
	register(new(CookieValueRequestCondition))
	register(new(UrlParamsRequestCondition))
	register(new(FormRequestCondition))
	register(new(CompositeRequestCondition))
	c := new(CompositeRequestCondition)
	c.Name = "composite-all"
	c.Add(new(HeaderValueRequestCondition))
	c.Add(new(CookieValueRequestCondition))
	c.Add(new(UrlParamsRequestCondition))
	c.Add(new(FormRequestCondition))
	register(c)

}
func register(rc RequestCondition) {

	RequestConditions[rc.Id()] = rc
}

type RequestCondition interface {
	Matched(ctx *RequestContext) bool
	Id() string
	Conf(conf kvs.ConfigSource)
}
type BaseConfRequestCondition struct {
	conf kvs.ConfigSource
}

func (b *BaseConfRequestCondition) Conf(conf kvs.ConfigSource) {
	b.conf = conf
}

//在conf中增加以下kv来满足：
// traffic.cond.AppName.client.addrs=192.168.1.2,192.168.1.3 多个IP逗号分割，
// traffic.cond.AppName.version=1.0.1 指定灰度版本来设定灰度版本
type ClientAddrRequestCondition struct {
	BaseConfRequestCondition
}

func (r *ClientAddrRequestCondition) Id() string {
	return "remoteAddr"
}
func (r *ClientAddrRequestCondition) Matched(ctx *RequestContext) bool {
	ip := getClientIP(ctx)
	ips := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.client.addrs", ctx.AppName), "")
	return strings.Contains(ips, ip)
}

func getClientIP(ctx *RequestContext) string {
	//获取客户端真实的IP地址
	raddr := ctx.Request.RemoteAddr
	xff := ctx.Request.Header.Get("X-Forwarded-For")
	if xff == "" {
		ra := ctx.Request.Header.Get("X-Real-IP")
		if ra != "" {
			raddr = ra
		}
	} else {
		xffs := strings.Split(xff, ",")
		raddr = xffs[0]
	}
	raddr = strings.Split(raddr, ":")[0]
	fmt.Println(raddr)
	return raddr

}

type HeaderValueRequestCondition struct {
	BaseConfRequestCondition
}

func (r *HeaderValueRequestCondition) Id() string {
	return "header"
}

func (r *HeaderValueRequestCondition) Matched(ctx *RequestContext) bool {
	key := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.header.key", ctx.AppName), "")
	values := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.header.values", ctx.AppName), "")

	v := ctx.Request.Header.Get(key)
	if strings.Contains(values, v) {
		return true
	}
	ranges := r.conf.Strings(fmt.Sprintf("traffic.cond.%s.header.ranges", ctx.AppName))
	return contains(ranges, v)
}

func contains(ranges []string, v string) bool {
	numVal, err := strconv.Atoi(v)
	if err != nil {
		return false
	}
	for _, r := range ranges {
		kv := kvs.NewKeyValueByStrDelims("", r, "-")
		nums := kv.Ints()
		if nums[0] <= numVal && numVal <= nums[0] {
			return true
		}
	}
	return false
}

type CookieValueRequestCondition struct {
	BaseConfRequestCondition
}

func (r *CookieValueRequestCondition) Id() string {
	return "cookie"
}
func (r *CookieValueRequestCondition) Matched(ctx *RequestContext) bool {
	key := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.cookie.key", ctx.AppName), "")
	values := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.cookie.values", ctx.AppName), "")
	c, err := ctx.Request.Cookie(key)
	if err != nil {
		return false
	}
	if strings.Contains(values, c.Value) {
		return true
	}
	ranges := r.conf.Strings(fmt.Sprintf("traffic.cond.%s.cookie.ranges", ctx.AppName))
	return contains(ranges, c.Value)
}

type UrlParamsRequestCondition struct {
	BaseConfRequestCondition
}

func (r *UrlParamsRequestCondition) Id() string {
	return "urlParams"
}
func (r *UrlParamsRequestCondition) Matched(ctx *RequestContext) bool {
	key := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.urlParams.key", ctx.AppName), "")
	values := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.urlParams.values", ctx.AppName), "")
	v := ctx.Request.URL.Query().Get(key)
	if strings.Contains(values, v) {
		return true
	}
	ranges := r.conf.Strings(fmt.Sprintf("traffic.cond.%s.urlParams.ranges", ctx.AppName))
	return contains(ranges, v)
}

type FormRequestCondition struct {
	BaseConfRequestCondition
}

func (r *FormRequestCondition) Id() string {
	return "form"
}
func (r *FormRequestCondition) Matched(ctx *RequestContext) bool {
	key := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.form.key", ctx.AppName), "")
	values := r.conf.GetDefault(fmt.Sprintf("traffic.cond.%s.form.values", ctx.AppName), "")
	v := ctx.Request.Form.Get(key)
	if strings.Contains(values, v) {
		return true
	}
	ranges := r.conf.Strings(fmt.Sprintf("traffic.cond.%s.form.ranges", ctx.AppName))
	return contains(ranges, v)
}

type CompositeRequestCondition struct {
	conf              kvs.ConfigSource
	RequestConditions []RequestCondition
	Name              string
}

func (b *CompositeRequestCondition) Conf(conf kvs.ConfigSource) {
	b.conf = conf
}
func (r *CompositeRequestCondition) Id() string {
	if r.Name == "" {
		return "composite"
	} else {
		return r.Name
	}
}

func (r *CompositeRequestCondition) Add(cond RequestCondition) {
	r.RequestConditions = append(r.RequestConditions, cond)
}

func (r *CompositeRequestCondition) Matched(ctx *RequestContext) bool {
	for _, cond := range r.RequestConditions {
		cond.Conf(r.conf)
		if cond.Matched(ctx) {
			return true
		}
	}
	return false
}
