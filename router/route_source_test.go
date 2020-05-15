package router

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestNewKeyValueRouteSource(t *testing.T) {

	Convey("1", t, func() {

		Convey("add", func() {
			p := NewKeyValueRouteSource()
			p.Init()
			r1 := &Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false}
			p.Add(r1)
			r2 := &Route{Source: "/app2/v1/user", ServiceId: "app2", Target: "/v1/user", StripPrefix: false}
			p.Add(r2)
			So(p.routes, ShouldNotBeNil)
			So(len(p.routes), ShouldEqual, 2)
			fmt.Println(p.Name())
		})
		Convey("update", func() {
			p := NewKeyValueRouteSource()
			p.Init()
			r1 := &Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false}
			p.Add(r1)
			r2 := &Route{Source: "/app2/v1/user", ServiceId: "app2", Target: "/v1/user", StripPrefix: false}
			p.Add(r2)
			So(p.routes, ShouldNotBeNil)
			//重复添加
			p.Add(r1)
			p.Add(r2)
			r3 := &Route{Source: "/app3/v1/user", ServiceId: "app3", Target: "/v1/user", StripPrefix: false}
			p.Add(r3)
			So(len(p.routes), ShouldEqual, 3)
		})

		Convey("In time sync", func() {
			p := NewKeyValueRouteSource()
			p.Init()
			router := &Router{}
			p.SetRouterChangedCallback(func(route *Route) {
				router.AddRoute(*route)
			})
			r1 := &Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false}
			r2 := &Route{Source: "/app2/v1/user", ServiceId: "app2", Target: "/v1/user", StripPrefix: false}
			So(p.routes, ShouldNotBeNil)
			//实时同步
			p.AddInTime(r1)
			p.AddInTime(r2)
			So(len(p.routes), ShouldEqual, 2)
			So(len(router.Routes), ShouldEqual, 2)
		})

		Convey("In time sync for update", func() {
			p := NewKeyValueRouteSource()
			p.Init()
			router := &Router{}
			p.SetRouterChangedCallback(func(route *Route) {
				router.AddRoute(*route)
			})
			r1 := &Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false}
			p.Add(r1)
			r2 := &Route{Source: "/app2/v1/user", ServiceId: "app2", Target: "/v1/user", StripPrefix: false}
			p.Add(r2)
			So(p.routes, ShouldNotBeNil)
			//实时同步
			p.AddInTime(r1)
			p.AddInTime(r2)
			So(len(p.routes), ShouldEqual, 2)
			So(len(router.Routes), ShouldEqual, 2)
		})

	})

}
