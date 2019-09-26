package main

import "time"

//DEFAULT
const (
//key
//default value
)

//bootstrap
const (
    //key
    KEY_BOOTSTRAP_ZK_ENABLED     = "bootstrap.zk.enabled"
    KEY_BOOTSTRAP_ZK_URLS        = "bootstrap.zk.urls"
    KEY_BOOTSTRAP_ZK_TIMEOUT     = "bootstrap.zk.timeout"
    KEY_BOOTSTRAP_ZK_ROOT        = "bootstrap.zk.root"
    KEY_BOOTSTRAP_CONSUL_ENABLED = "bootstrap.consul.enabled"
    KEY_BOOTSTRAP_CONSUL_ADDRESS = "bootstrap.consul.address"
    KEY_BOOTSTRAP_CONSUL_ROOT    = "bootstrap.consul.root"
    //default value
    DEFAULT_BOOTSTRAP_ZK_ENABLED     = "false"
    DEFAULT_BOOTSTRAP_ZK_URLS        = "127.0.0.1:2181"
    DEFAULT_BOOTSTRAP_ZK_TIMEOUT     = time.Duration(10000) * time.Millisecond
    DEFAULT_BOOTSTRAP_ZK_ROOT        = "/zebra/bootstrap"
    DEFAULT_BOOTSTRAP_CONSUL_ENABLED = "false"
    DEFAULT_BOOTSTRAP_CONSUL_ADDRESS = "127.0.0.1:8500"
    DEFAULT_BOOTSTRAP_CONSUL_ROOT    = "zebra/bootstrap"
)

//server
const (
    //key
    KEY_SERVER_DEBUG            = "server.debug"
    KEY_SERVER_PORT             = "server.port"
    KEY_SERVER_MODE             = "server.mode"
    KEY_SERVER_FAVICON_ICO_PATH = "server.favicon.ico.path"
    KEY_SERVER_GMS_ENABLED      = "server.gms.enabled"
    KEY_SERVER_GMS_DOMAIN       = "server.gms.domain"
    KEY_SERVER_GMS_PORT         = "server.gms.port"
    KEY_SERVER_GMS_DEBUG        = "server.gms.debug"
    //default value
    DEFAULT_SERVER_DEBUG            = "true"
    DEFAULT_SERVER_PORT             = 19001
    DEFAULT_SERVER_MODE             = "client"
    DEFAULT_SERVER_FAVICON_ICO_PATH = "favicon.ico"
    DEFAULT_SERVER_GMS_ENABLED      = "true"
    DEFAULT_SERVER_GMS_DOMAIN       = ""
    DEFAULT_SERVER_GMS_PORT         = 17980
    DEFAULT_SERVER_GMS_DEBUG        = "true"
)

//hystrix
const (
    //key
    KEY_HYSTRIX_CIRCUIT_ENABLED                = "hystrix.circuit.enabled"
    KEY_HYSTRIX_DEFAULT_TIMEOUT                = "hystrix.default.Timeout"
    KEY_HYSTRIX_DEFAULT_MAXCONCURRENTREQUESTS  = "hystrix.default.MaxConcurrentRequests"
    KEY_HYSTRIX_DEFAULT_REQUESTVOLUMETHRESHOLD = "hystrix.default.RequestVolumeThreshold"
    KEY_HYSTRIX_DEFAULT_SLEEPWINDOW            = "hystrix.default.SleepWindow"
    KEY_HYSTRIX_DEFAULT_ERRORPERCENTTHRESHOLD  = "hystrix.default.ErrorPercentThreshold"
    //default value
    DEFAULT_HYSTRIX_CIRCUIT_ENABLED                = "true"
    DEFAULT_HYSTRIX_DEFAULT_TIMEOUT                = 6000
    DEFAULT_HYSTRIX_DEFAULT_MAXCONCURRENTREQUESTS  = 100
    DEFAULT_HYSTRIX_DEFAULT_REQUESTVOLUMETHRESHOLD = 20
    DEFAULT_HYSTRIX_DEFAULT_SLEEPWINDOW            = 5000
    DEFAULT_HYSTRIX_DEFAULT_ERRORPERCENTTHRESHOLD  = 50
)

//http
const (
    //key
    KEY_HTTP_SERVER_READTIMEOUT     = "http.server.ReadTimeout"
    KEY_HTTP_SERVER_WRITETIMEOUT    = "http.server.WriteTimeout"
    KEY_HTTP_CLIENT_CONNECT_TIMEOUT = "http.client.connect.timeout"
    //default value
    DEFAULT_HTTP_SERVER_READTIMEOUT     = 10000
    DEFAULT_HTTP_SERVER_WRITETIMEOUT    = 10000
    DEFAULT_HTTP_CLIENT_CONNECT_TIMEOUT = 3000
)

//lb
const (
    //key
    KEY_LB_DEFAULT_BALANCER_NAME    = "lb.default.balancer.name"
    KEY_LB_DEFAULT_MAX_FAILS        = "lb.default.max.fails"
    KEY_LB_DEFAULT_FAIL_TIME_WINDOW = "lb.default.fail.time.window"
    KEY_LB_DEFAULT_FAIL_SLEEP_MODE  = "lb.default.fail.sleep.mode"
    KEY_LB_DEFAULT_FAIL_SLEEP_X     = "lb.default.fail.sleep.x"
    KEY_LB_DEFAULT_FAIL_SLEEP_MAX   = "lb.default.fail.sleep.max"
    KEY_LB_LOG_SELECTED_ENABLED     = "lb.log.selected.enabled"
    KEY_LB_LOG_SELECTED_ALL         = "lb.log.selected.all"
    //default value
    DEFAULT_LB_DEFAULT_BALANCER_NAME    = "WeightRobinRound"
    DEFAULT_LB_DEFAULT_MAX_FAILS        = 3
    DEFAULT_LB_DEFAULT_FAIL_TIME_WINDOW = time.Duration(10000) * time.Millisecond
    DEFAULT_LB_DEFAULT_FAIL_SLEEP_MODE  = "seq"
    DEFAULT_LB_DEFAULT_FAIL_SLEEP_X     = "1,1,2,3,5,8,13,21"
    DEFAULT_LB_DEFAULT_FAIL_SLEEP_MAX   = time.Duration(60000) * time.Millisecond
    DEFAULT_LB_LOG_SELECTED_ENABLED     = "true"
    DEFAULT_LB_LOG_SELECTED_ALL         = "false"
)

//routes
const (
    //key
    KEY_ROUTES_STRIP_PREFIX = "routes.strip.prefix"
    //default value
    DEFAULT_ROUTES_STRIP_PREFIX = "false"
)

//ini
const (
    //key
    KEY_INI_ROUTES_ENABLED    = "ini.routes.enabled"
    KEY_INI_DISCOVERY_ENBALED = "ini.discovery.enbaled"
    KEY_INI_DISCOVERY_DIR     = "ini.discovery.dir"
    //default value
    DEFAULT_INI_ROUTES_ENABLED    = "true"
    DEFAULT_INI_DISCOVERY_ENBALED = "false"
    DEFAULT_INI_DISCOVERY_DIR     = "services"
)

//eureka
const (
    //key
    KEY_EUREKA_SERVER_ENABLED     = "eureka.server.enabled"
    KEY_EUREKA_SERVER_URLS        = "eureka.server.urls"
    KEY_EUREKA_DISCOVERY_INTERVAL = "eureka.discovery.interval"
    //default value
    DEFAULT_EUREKA_SERVER_ENABLED     = "false"
    DEFAULT_EUREKA_SERVER_URLS        = "http://127.0.0.1:8761/eureka"
    DEFAULT_EUREKA_DISCOVERY_INTERVAL = time.Duration(10000) * time.Millisecond
)

//consul
const (
    //key
    KEY_CONSUL_ENABLED            = "consul.enabled"
    KEY_CONSUL_ADDRESS            = "consul.address"
    KEY_CONSUL_ROUTES_ENABLED     = "consul.routes.enabled"
    KEY_CONSUL_ROUTES_ROOT        = "consul.routes.root"
    KEY_CONSUL_DISCOVERY_ENABLED  = "consul.discovery.enabled"
    KEY_CONSUL_DISCOVERY_ADDRESS  = "consul.discovery.address"
    KEY_CONSUL_DISCOVERY_INTERVAL = "consul.discovery.interval"
    //default value
    DEFAULT_CONSUL_ENABLED            = "false"
    DEFAULT_CONSUL_ADDRESS            = "127.0.0.1:8500"
    DEFAULT_CONSUL_ROUTES_ENABLED     = "true"
    DEFAULT_CONSUL_ROUTES_ROOT        = "zebra/routes"
    DEFAULT_CONSUL_DISCOVERY_ENABLED  = "true"
    DEFAULT_CONSUL_DISCOVERY_ADDRESS  = "${consul.address}"
    DEFAULT_CONSUL_DISCOVERY_INTERVAL = time.Duration(10000) * time.Millisecond
)

//zk
const (
    //key
    KEY_ZK_ENABLED           = "zk.enabled"
    KEY_ZK_CONN_URLS         = "zk.conn.urls"
    KEY_ZK_TIMEOUT           = "zk.timeout"
    KEY_ZK_ROOT              = "zk.root"
    KEY_ZK_ROUTES_CONTEXTS   = "zk.routes.contexts"
    KEY_ZK_ROUTES_WATCH_TYPE = "zk.routes.watch.type"
    KEY_ZK_DISCOVERY_ENBALED = "zk.discovery.enbaled"
    KEY_ZK_DISCOVERY_PATH    = "zk.discovery.path"
    //default value
    DEFAULT_ZK_ENABLED           = "true"
    DEFAULT_ZK_CONN_URLS         = "127.0.0.1:2181"
    DEFAULT_ZK_TIMEOUT           = time.Duration(3000) * time.Millisecond
    DEFAULT_ZK_ROOT              = "/zebra"
    DEFAULT_ZK_ROUTES_CONTEXTS   = "routes"
    DEFAULT_ZK_ROUTES_WATCH_TYPE = "all"
    DEFAULT_ZK_DISCOVERY_ENBALED = "false"
    DEFAULT_ZK_DISCOVERY_PATH    = "/services/"
)

//sql
const (
    //key
    KEY_SQL_ROUTES_ENABLED     = "sql.routes.enabled"
    KEY_SQL_DRIVERNAME         = "sql.driverName"
    KEY_SQL_URL                = "sql.url"
    KEY_SQL_DISCOVERY_ENABLED  = "sql.discovery.enabled"
    KEY_SQL_DISCOVERY_INTERVAL = "sql.discovery.interval"
    //default value
    DEFAULT_SQL_ROUTES_ENABLED     = "true"
    DEFAULT_SQL_DRIVERNAME         = "mysql"
    DEFAULT_SQL_URL                = "root:kry02Local@?DB@tcp(172.16.1.248:3306)/po?charset=utf8"
    DEFAULT_SQL_DISCOVERY_ENABLED  = "false"
    DEFAULT_SQL_DISCOVERY_INTERVAL = time.Duration(10000) * time.Millisecond
)

//k8s
const (
    //key
    KEY_K8S_ENABLED            = "k8s.enabled"
    KEY_K8S_URLS               = "k8s.urls"
    KEY_K8S_DISCOVERY_INTERVAL = "k8s.discovery.interval"
    //default value
    DEFAULT_K8S_ENABLED            = "false"
    DEFAULT_K8S_URLS               = "http://127.0.0.1:8080"
    DEFAULT_K8S_DISCOVERY_INTERVAL = time.Duration(10000) * time.Millisecond
)

//metrics
const (
    //key
    KEY_METRICS_EXPORT_INFLUX_ENABLED  = "metrics.export.influx.enabled"
    KEY_METRICS_EXPORT_INFLUX_INTERVAL = "metrics.export.influx.interval"
    KEY_METRICS_EXPORT_INFLUX_URL      = "metrics.export.influx.url"
    KEY_METRICS_EXPORT_INFLUX_DATABASE = "metrics.export.influx.database"
    KEY_METRICS_EXPORT_INFLUX_USERNAME = "metrics.export.influx.username"
    KEY_METRICS_EXPORT_INFLUX_PASSWORD = "metrics.export.influx.password"
    //default value
    DEFAULT_METRICS_EXPORT_INFLUX_ENABLED  = "false"
    DEFAULT_METRICS_EXPORT_INFLUX_INTERVAL = time.Duration(10000) * time.Millisecond
    DEFAULT_METRICS_EXPORT_INFLUX_URL      = "http://127.0.0.1:8086"
    DEFAULT_METRICS_EXPORT_INFLUX_DATABASE = "metrics"
    DEFAULT_METRICS_EXPORT_INFLUX_USERNAME = ""
    DEFAULT_METRICS_EXPORT_INFLUX_PASSWORD = ""
)
