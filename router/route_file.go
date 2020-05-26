package router

import (
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/go-utils"
	pini "github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	"github.com/tietang/zebra/health"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime/debug"
)

const (
	KEY_INI_ROUTES_ENABLED          = "ini.routes.enabled"
	KEY_INI_ROUTES_DIR              = "ini.routes.dir"
	DEFAULT_INI_ROUTES_RELATIVE_DIR = "routes"

	KEY_INI_ROUTES_STRIP_PREFIX = "ini.routes.strip.prefix"
)

/*

[app1]
/app1/v1/user=/v1/user
/app1/info=/info
/api/info=/info,false


*/
///api/v1/user,user,/api/v1/user,false
// format: source pattern = service name[, target pattern][,stripPrefix]
type IniFileRouteSource struct {
	KeyValueRouteSource
	routesDir string
}

func NewIniFileRouteSource(conf kvs.ConfigSource) *IniFileRouteSource {
	dir := conf.GetDefault(KEY_INI_ROUTES_DIR, DEFAULT_INI_ROUTES_RELATIVE_DIR)
	routesDir := dir
	//默认没有配置kvs.routes.dir或者配置为相对路径，则在程序运行目录查找
	if !filepath.IsAbs(routesDir) {
		pwd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		routesDir = filepath.Join(pwd, dir)
	}

	isExists, err := utils.PathExists(routesDir)
	//如果不存在,则在配置文件目录查找
	if !isExists || err != nil {
		log.Warn("no exists routes running dir : ", routesDir, ", err:", err)
		currDir, err := conf.Get(pini.KEY_INI_CURRENT_DIR)
		//未通过props来启动server，则不存在该配置
		if err != nil || currDir == "" {
			log.Warn("no props dir: ", currDir, err)
		} else {
			routesDir = filepath.Join(currDir, DEFAULT_INI_ROUTES_RELATIVE_DIR)
		}

	}
	log.Info("routes dir: ", routesDir)
	k := &IniFileRouteSource{
		routesDir: routesDir,
	}
	k.isEnabled = true
	isExists, err = utils.PathExists(routesDir)
	if !isExists || err != nil {
		k.isEnabled = false
		log.Warn("not exists dir:", routesDir)
	}
	k.conf = conf
	k.StripPrefix = conf.GetBoolDefault(KEY_INI_ROUTES_STRIP_PREFIX, false)
	k.name = "dir: " + dir
	k.routes = make([]*Route, 0)
	return k
}

func (p *IniFileRouteSource) Init() {
	if p.isInit {
		return
	}
	if p.isEnabled {
		p.readDir()
	}
	p.isInit = true
}
func (p *IniFileRouteSource) readDir() {
	if p.routesDir == "" {
		log.Warn("routes path is empty ")
		return
	}

	log.Info("load routes by path: ", p.routesDir)
	dir, err := ioutil.ReadDir(p.routesDir)
	if err == nil {
		for _, fi := range dir {
			if fi.IsDir() { // 忽略目录
				continue
			}
			file := filepath.Join(p.routesDir, fi.Name())
			log.Info("read routes file: ", file)
			p.readIniFile(file)
		}
	} else {
		log.Error(string(debug.Stack()))
		log.Warn(err, ", path=", p.routesDir)
	}
}

func (p *IniFileRouteSource) readIniFile(file string) {
	log.Info("load routes file: ", file)
	iniFile, err := ini.Load(file)
	if err != nil {
		log.Warn(err)
		return
	}
	sections := iniFile.Sections()
	sp := ReadIniSections(p.routerStripPrefix, sections, func(route *Route) {
		p.Add(route)
	})
	p.addConfigSource(sp)

}

func (h *IniFileRouteSource) readHosts() {

}
func (h *IniFileRouteSource) CheckHealth(rootHealth *health.RootHealth) {
	_, err := os.Stat(h.routesDir)
	ok, desc := true, "ok"
	if err == nil {
		ok, desc = true, "ok"
	}
	if os.IsNotExist(err) {
		ok, desc = false, "Not Exists "+h.name
	}

	if ok {
		rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_UP, Desc: desc}

	} else {
		rootHealth.Healths["eureka server"] = &health.Health{Status: health.STATUS_DOWN, Desc: desc}
		rootHealth.Status = health.STATUS_DOWN
	}
}
