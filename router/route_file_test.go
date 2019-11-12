package router

import (
    "fmt"
    "github.com/tietang/props/kvs"
    "os"
    "path/filepath"
    "testing"
)

func TestNewIniRouteSource(t *testing.T) {
    dir, _ := os.Getwd()
    dir = filepath.Join(dir, "testdata")
    conf := kvs.NewEmptyMapConfigSource("test_map")
    conf.Set("ini.routes.dir", "testdata")
    //conf.Set("kvs.routes.dir","")
    p := NewIniFileRouteSource(conf)
    p.Init()
    Convey("1", t, func() {

        So(p.routes, ShouldNotBeNil)
        So(len(p.routes), ShouldEqual, 2)
        fmt.Println(p.Name())
    })

}
