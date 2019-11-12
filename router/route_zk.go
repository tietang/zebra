package router

import (
    "bytes"
    "github.com/samuel/go-zookeeper/zk"
    log "github.com/sirupsen/logrus"
    "github.com/tietang/props/ini"
    "github.com/tietang/props/kvs"
    pzk "github.com/tietang/props/zk"
    "github.com/tietang/zebra/health"
    "path"
    "strings"
    "time"
)

const (
    KEY_ZK_ENABLED           = "zk.enabled"
    KEY_ZK_CONN_URLS         = "zk.conn.urls"
    KEY_ZK_CONN_TIMEOUT      = "zk.conn.timeout"
    KEY_ZK_ROOT              = "zk.root"
    KEY_ZK_ROUTES_CONTEXTS   = "zk.routes.contexts"
    KEY_ZK_ROUTES_WATCH_TYPE = "zk.routes.watch.type"

    //
    DATA_KEY_ROUTES_NODE   = "routes"
    KEY_ROUTES_NODE_PREFIX = "routes"
    KEY_ROUTES_NOTIFY_NODE = "notify"
    KEY_ROUTE_STRIP_PREFIX = "routes.strip.prefix"

    DEFAULT_ZK_CONN_TIMEOUT = 3 * time.Second
    DEFAULT_ZK_ROOT         = "/zebra/root"
    DEFAULT_ROUTES_CONTEXT  = "routes"

    //
    WATCH_TYPE_ALL         = "ALL"
    WATCH_TYPE_NOTIFY_NODE = "NOTIFY-NODE"
)

type ZookeeperRouteSource struct {
    KeyValueRouteSource
    ZkConfigSource    *kvs.CompositeConfigSource
    zkConn            *zk.Conn
    zkConnUrls        []string
    zkRootPath        string
    zkContexts        []string
    zkRoutesWatchType string

    globalStripPrefix bool
}

func NewZookeeperRouteSource(conf kvs.ConfigSource) *ZookeeperRouteSource {
    z := &ZookeeperRouteSource{}
    z.conf = conf

    z.routes = make([]*Route, 0)
    globalStripPrefix := conf.GetBoolDefault(KEY_ROUTE_STRIP_PREFIX, false)
    urls := conf.GetDefault(KEY_ZK_CONN_URLS, "127.0.0.1:2181")
    watchType := conf.GetDefault(KEY_ZK_ROUTES_WATCH_TYPE, "all")
    contextsStr := conf.GetDefault(KEY_ZK_ROUTES_CONTEXTS, DEFAULT_ROUTES_CONTEXT)
    root := conf.GetDefault(KEY_ZK_ROOT, DEFAULT_ZK_ROOT)
    contexts := strings.Split(contextsStr, ",")
    zkUrls := strings.Split(urls, ",")
    //
    z.zkRootPath = root
    z.zkRoutesWatchType = watchType
    z.name = "zookeeper:" + root
    z.globalStripPrefix = globalStripPrefix
    //
    ctx := make([]string, 0)
    idx := strings.LastIndex(root, "/")
    if idx == -1 {
        root = root + "/"
    }
    for _, c := range contexts {
        ctx = append(ctx, path.Join(root, c))
    }
    z.zkContexts = ctx
    z.zkConnUrls = zkUrls

    return z
}

func (z *ZookeeperRouteSource) Init() {
    if z.isInit {
        return
    }
    timeout := z.conf.GetDurationDefault(KEY_ZK_CONN_TIMEOUT, DEFAULT_ZK_CONN_TIMEOUT)
    conn, ch, err := zk.Connect(z.zkConnUrls, timeout)

    if err != nil {
        log.Error(err)
        panic(err)
    }
    //waiting for connected.
    for {
        event := <-ch
        log.Info(event)
        if event.State == zk.StateConnected {
            break
        }
    }

    z.zkConn = conn
    source := pzk.NewZookeeperCompositeConfigSource(z.zkContexts, z.zkConnUrls, timeout)
    //source := pzk.NewZookeeperCompositeConfigSourceByConn(z.zkContexts, conn)
    z.ZkConfigSource = source
    z.loadRouteByConfigSource(z.ZkConfigSource)
    z.WatchRoutesPath()
}

func (g *ZookeeperRouteSource) WatchRoutesPath() {

    for _, zkPath := range g.zkContexts {
        if strings.ToUpper(g.zkRoutesWatchType) == WATCH_TYPE_ALL {

            paths, _, err := g.zkConn.Children(zkPath)
            if err == nil {
                for _, child := range paths {
                    if !strings.HasSuffix(child, KEY_ROUTES_NOTIFY_NODE) && !strings.HasSuffix(child, "foo") {
                        go g.watchGet(g.zkConn, path.Join(zkPath, child))
                    }
                }
            } else {
                log.Error(err)
            }
        } else {
            go g.watchGet(g.zkConn, path.Join(zkPath, KEY_ROUTES_NOTIFY_NODE))
        }
    }
}

func (g *ZookeeperRouteSource) watchChildren(zkConn *zk.Conn, pathStr string) {
    children, stat, ch, err := zkConn.ChildrenW(pathStr)
    if err != nil {
        //panic(err)
        log.Error(err, "can't watch children for parent node:", pathStr)
    }
    log.Infof("%+v %+v\n", children, stat)
    e := <-ch
    log.Infof("%+v\n", e)
    g.watchChildren(zkConn, pathStr)
}

func (g *ZookeeperRouteSource) watchGet(zkConn *zk.Conn, pathStr string) {
    log.Info("watched path: ", pathStr)
    exists, _, _ := zkConn.Exists(pathStr)
    if !exists {
        zkConn.Create(pathStr, []byte("1"), 1, nil)
    }
    children, stat, ch, err := zkConn.GetW(pathStr)
    if err != nil {
        //panic(err)
        log.Error(err, "can't watch node:", pathStr)
    }
    //log.Info("watched children: ", string(children), stat)
    log.Info("watched children: ", pathStr, stat)
    e := <-ch

    //pPath:=path.Dir(e.Path)
    //source := kvs.NewZookeeperCompositeConfigSourceByConn(g.zkContexts, g.zkConn)
    r := bytes.NewReader(children)
    source := ini.NewIniFileConfigSourceByReader("string-reader", r)
    g.loadRouteByConfigSource(source)
    //TODO reload router

    log.Info("notify: ", e)
    g.watchGet(zkConn, pathStr)
}

func (g *ZookeeperRouteSource) CheckHealth(rootHealth *health.RootHealth) {
    if g.zkRootPath != "" && g.zkConn != nil {
        s := g.zkConn.State()
        if s == zk.StateConnected {
            rootHealth.Healths["zookeeper"] = &health.Health{Status: health.STATUS_UP}
        } else {
            rootHealth.Healths["zookeeper"] = &health.Health{Status: health.STATUS_DOWN}
            rootHealth.Status = health.STATUS_DOWN
        }
    }
}
