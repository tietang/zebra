package router

import (
	"fmt"
	"github.com/go-ini/ini"
	log "github.com/sirupsen/logrus"
	"github.com/tietang/props/kvs"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

const (
	DISCOVERY_INTERVAL_DEFAULT        = 10 * time.Second
	KEY_LABEL_ROUTE_PREFIX            = "route_prefix"
	DEFAULT_INI_ROUTES_SERVICE_SUFFIX = ".routes"
)

type Route struct {
	Id                 string
	ServiceId          string
	Source             string
	SourcePrefix       string
	SourceIsFuzzyMatch bool
	Target             string
	TargetIsFuzzyMatch bool
	TargetPrefix       string
	StripPrefix        bool
	//
	Subdomain string // "admin."
	Path      string // "/api/user/:id"
	//DataCenter string
	ServiceSource string
	IsForceUpdate bool
}

//幂等
func (r *Route) Init() {
	sourceFuzzyMatchIndex := strings.LastIndex(r.Source, "/**")
	if sourceFuzzyMatchIndex > -1 {
		r.SourcePrefix = strings.TrimSuffix(r.Source, "/**")
		r.SourceIsFuzzyMatch = true
	}
	targetFuzzyMatchIndex := strings.LastIndex(r.Target, "/**")
	if targetFuzzyMatchIndex > -1 {
		r.TargetPrefix = strings.TrimSuffix(r.Target, "/**")
		r.TargetIsFuzzyMatch = true
	}
	r.Path = r.Source
	r.Id = strings.Join([]string{r.ServiceId, r.Source}, ":")
}

func (route *Route) GetRouteTargetPath(sourcePath string) string {

	tpath := sourcePath
	isStrip := false

	if route.Target != "" {

		if route.TargetIsFuzzyMatch {
			tpath = sourcePath
			isStrip = true
		} else {
			return route.Target
		}
	} else {
		isStrip = true
	}

	if isStrip && route.StripPrefix {

		//print(strings.format("%s %d %d",tpath,index,strings.len(host.prefix)))
		if route.TargetIsFuzzyMatch {
			return route.TargetPrefix + strings.TrimPrefix(tpath, route.SourcePrefix)
		} else {
			return strings.TrimPrefix(tpath, route.SourcePrefix)
		}
	} else {
		if route.TargetIsFuzzyMatch {
			return route.TargetPrefix + tpath
		} else {
			return tpath
		}
	}
}

func (route *Route) isMatch(sourcePath string) bool {
	if route.SourceIsFuzzyMatch {
		hasPrefix := strings.HasPrefix(sourcePath, route.SourcePrefix)
		if hasPrefix {
			return true
		}
	} else {
		if sourcePath == route.Source {
			return true
		}
	}
	return false
}

// URL creates a URL using the current host and the given parameters.
// The parameters should be given in the sequence of name1, value1, name2, value2, and so on.
// If a parameter in the host is not provided a value, the parameter token will remain in the resulting URL.
// The method will perform URL encoding for all given parameter values.
func (r *Route) URL(pairs ...interface{}) (s string) {
	for i := 0; i < len(pairs); i++ {
		name := fmt.Sprintf("%s/**", pairs[i])
		value := ""
		if i < len(pairs)-1 {
			value = url.QueryEscape(fmt.Sprint(pairs[i+1]))
		}
		s = strings.Replace(s, name, value, -1)
	}
	return
}

//ini格式
func parseRouteLine(routerStripPrefix bool, appName, key, value string) *Route {

	source := strings.TrimSpace(key)
	app := strings.TrimSpace(appName)
	routeStr := strings.TrimSpace(value)
	vs := strings.Split(routeStr, ",")
	size := len(vs)
	if size < 1 {
		return nil
	}
	stripPrefix := false

	var target string
	var err error

	// /api/v1/user=/api/v1/user
	if size >= 1 {
		target = vs[0]
	}
	// /api/v1/user=/api/v1/user,false
	if size >= 2 {
		stripPrefix, err = strconv.ParseBool(vs[1])
		if err != nil {
			stripPrefix = routerStripPrefix
		}
	}

	route := &Route{
		Source:      source,
		ServiceId:   app,
		Target:      target,
		StripPrefix: stripPrefix,
	}
	log.Debug("add router by config: ", routeStr)

	return route
}

func ReadIniSections(routerStripPrefix bool, sections []*ini.Section, call func(route *Route)) kvs.ConfigSource {
	sp := kvs.NewEmptyMapConfigSource("services.props")
	for _, section := range sections {
		name := section.Name()
		if strings.HasSuffix(name, DEFAULT_INI_ROUTES_SERVICE_SUFFIX) {
			//s[len(s)-len(suffix):] == suffix
			idx := len(name) - len(DEFAULT_INI_ROUTES_SERVICE_SUFFIX)
			serviceId := name[:idx]
			for _, kv := range section.Keys() {
				key := kv.Name()
				value := kv.String()
				route := parseRouteLine(routerStripPrefix, serviceId, key, value)
				call(route)
				//log.WithField("host", host).Info("add host: ")
				log.Debug("add host: ", key, "=", value)
			}
		} else {
			serviceId := name
			for _, kv := range section.Keys() {
				key := serviceId + "." + kv.Name()
				value := kv.String()
				sp.Set(key, value)
			}
		}
	}

	return sp

}

func toPath(str string) string {
	if str == "" {
		return ""
	}
	urlPatterns := strings.Split(str, ".")
	for i, v := range urlPatterns {
		if i == 0 && v == "" {
			urlPatterns[i] = "/"
		}
		if i == len(urlPatterns)-1 && v == "" {
			urlPatterns[i] = "**"
		}
	}
	urlPattern := path.Join(urlPatterns...)
	if strings.Index(urlPattern, "/") > 0 {
		urlPattern = path.Join("/", urlPattern)
	}

	return urlPattern
}
