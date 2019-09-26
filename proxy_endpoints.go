package zebra

import (
	"encoding/json"
	"github.com/fukata/golang-stats-api-handler"
	"github.com/rcrowley/go-metrics"
	"github.com/thoas/stats"
	"github.com/tietang/hystrix-go/hystrix"
	"github.com/tietang/zebra/meter"
	"github.com/tietang/zebra/router"
	"github.com/tietang/zebra/utils"
	"time"
)

var startTime string
var faviconIconData []byte

func (h *HttpProxyServer) initEndpoints() {
	h.rootEndpoint()
	h.healthEndpoint()
	h.infoEndpoint()
	h.metricsEndpoint()
	h.hystrixStreamEndpoint()
	//h.statsEndPoint()
	h.stats0Endpoint()
	h.gmsEndpoint()
	h.routesEndpoint()
	h.hostsEndpoint()
	h.faviconIconEndpoint()
	//dog
	//h.dogStatsEndpoint()
	//h.dogSysEndpoint()

}

func (h *HttpProxyServer) rootEndpoint() {
	s := struct {
		Version string
		Name    string
	}{
		Version: Version,
		Name:    Name,
	}
	h.Get("/", func(ctx *Context) error {

		data, err := json.Marshal(s)
		if err != nil {
			return err
		}
		ctx.Write(data)
		return nil
	})

}

//func (h *HttpProxyServer) dogSysEndpoint() {
//    h.Get("/dog/sys", func(ctx *Context) error {
//        dog.SigarHttpHandler()(ctx.ResponseWriter, ctx.Request)
//        return nil
//    })
//}
//
//func (h *HttpProxyServer) dogStatsEndpoint() {
//    h.Get("/dog/stats", func(ctx *Context) error {
//        stats_api.Handler(ctx.ResponseWriter, ctx.Request)
//        return nil
//    })
//}

func (h *HttpProxyServer) hostsEndpoint() {
	h.Get("/hosts", func(ctx *Context) error {

		kv1 := make(map[string]interface{})
		router.DefaultHosts.Range(func(key, value interface{}) bool {
			kv1[key.(string)] = value
			return true
		})
		kv2 := make(map[string]interface{})
		router.DefaultUnavailableHosts.Range(func(key, value interface{}) bool {
			kv2[key.(string)] = value
			return true
		})

		kvs := make(map[string]map[string]interface{})
		kvs["hosts"] = kv1
		kvs["unavailableHosts"] = kv2
		data, err := json.Marshal(kvs)
		if err != nil {
			return err
		}
		ctx.Write(data)
		return nil
	})

}
func (h *HttpProxyServer) faviconIconEndpoint() {
	h.Get("/favicon.ico", func(ctx *Context) error {
		//if faviconIconData == nil || len(faviconIconData) == 0 {
		//    path := h.conf.GetDefault(KEY_FAVICON_ICO_PATH, "favicon.ico")
		//    data, err := ioutil.ReadFile(path)
		//    if err != nil {
		//        return err
		//    }
		//    faviconIconData = data
		//}

		ctx.SetContentType("image/x-icon")
		//ctx.Write(faviconIconData)
		ctx.Write(utils.ICON_DATA)
		return nil
	})

}

func (h *HttpProxyServer) routesEndpoint() {
	h.Get("/routes", func(ctx *Context) error {

		data, err := json.Marshal(router.DefaultRouters)
		if err != nil {
			return err
		}
		ctx.Write(data)
		return nil
	})

}
func (h *HttpProxyServer) healthEndpoint() {

	h.Get("/health", func(ctx *Context) error {
		h.health.Check()
		data, err := json.Marshal(h.health)
		if err != nil {
			return err
		}
		ctx.Write(data)
		return nil
	})

}

func (h *HttpProxyServer) metricsEndpoint() {
	registries := map[string]metrics.Registry{
		"metricsDefaults": metrics.DefaultRegistry,
		"meterDefaults":   meter.DefaultRegistry,
		"urls":            router.UrlRegistry,
		"services":        router.ServiceRegistry,
		"instances":       router.InstanceRegistry,
	}
	h.Get("/metrics", func(ctx *Context) error {

		data, err := meter.MarshalJSON(registries)

		if err != nil {
			return err
		}
		ctx.Write(data)
		return nil
	})
}

func (h *HttpProxyServer) infoEndpoint() {
	startTime = time.Now().Format("2006-01-02T15:04:05.999999-07:00")
	h.Get("/info", func(ctx *Context) error {
		m := make(map[string]string)
		m["startTime"] = startTime
		data, err := json.Marshal(m)
		if err != nil {
			return err
		}
		ctx.Write(data)
		return nil
	})
}

func (h *HttpProxyServer) hystrixStreamEndpoint() {
	hs := hystrix.NewStreamHandler()
	hs.Start()
	h.Get("/hystrix.stream", func(ctx *Context) error {
		hs.ServeHTTP(ctx.ResponseWriter, ctx.Request)
		return nil
	})
}

func (h *HttpProxyServer) stats0Endpoint() {

	h.Get("/stats0", func(ctx *Context) error {
		stats_api.Handler(ctx.ResponseWriter, ctx.Request)
		return nil
	})

}

// Stats provides response time, status code count, etc.
var StatsMiddleware *stats.Stats

func (h *HttpProxyServer) statsEndPoint() {
	if StatsMiddleware == nil {
		StatsMiddleware = stats.New()
	}
	h.Use(func(ctx *Context) error {
		beginning, _ := StatsMiddleware.Begin(ctx.ResponseWriter)
		ctx.Next()
		StatsMiddleware.End(beginning)
		return nil
	})
	h.Get("/stats", func(ctx *Context) error {
		ctx.ResponseWriter.Header().Set("Content-Type", "application/json")

		stats := StatsMiddleware.Data()

		b, err := json.Marshal(stats)
		if err != nil {
			return err
		}
		ctx.Write(b)
		return nil
	})
}

func (p *HttpProxyServer) StartStatsServer() {
	conf := p.conf
	if conf.GetBoolDefault(KEY_SERVER_GMS_ENABLED, true) {

		c := &ClientConfig{
			Domain:           "",
			Port:             3009,
			PollInterval:     1000,
			Debug:            false,
			LogHostInfo:      true,
			LogCPUInfo:       true,
			LogTotalCPUTimes: true,
			LogPerCPUTimes:   true,
			LogMemory:        true,
			LogGoMemory:      true,
		}
		//config := &gms.ServerConfig{
		//    Domain: conf.GetDefault(KEY_GMS_DOMAIN, ""),
		//    Port:   conf.GetIntDefault(KEY_GMS_PORT, DEFAULT_GMS_PORT),
		//    Debug:  conf.GetBoolDefault(KEY_GMS_DEBUG, false),
		//}

		p.stats = new(Stats)

		if c.LogHostInfo {
			p.stats.GetHostInfo()
		}

		if c.LogCPUInfo {
			p.stats.GetCPUInfo()
		}
		//p.Use(func(ctx *Context) error {
		//    hr := NewHTTPRequest(ctx.ResponseWriter, ctx.Request)
		//    hr.HttpProxyServer = p
		//    err := ctx.Next()
		//    if err != nil {
		//        hr.Failure(err.Error())
		//    }
		//    hr.Complete()
		//    return nil
		//})
		ticker := time.NewTicker(time.Millisecond * time.Duration(c.PollInterval))
		defer ticker.Stop()

		for {
			select {
			case <-ticker.C:

				if c.LogTotalCPUTimes {
					p.stats.GetTotalCPUTimes()
				}

				if c.LogPerCPUTimes {
					p.stats.GetCPUTimes()
				}

				p.stats.GetMemoryInfo(c.LogMemory, c.LogGoMemory)

				//stats.HTTPRequests = c.httpStats.extract()
			}
		}
	}
}

func (h *HttpProxyServer) gmsEndpoint() {
	conf := h.conf
	isEnabled := conf.GetBoolDefault(KEY_SERVER_GMS_ENABLED, true)
	h.Get("/gms", func(ctx *Context) error {
		if isEnabled {
			data, err := json.Marshal(h.stats)
			if err != nil {
				return err
			}
			ctx.Write(data)
		} else {
			data, err := json.Marshal(&struct {
				Enabled bool
			}{Enabled: isEnabled})
			if err != nil {
				return err
			}
			ctx.Write(data)
		}
		return nil
	})

}
