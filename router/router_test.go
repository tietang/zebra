package router

import (
	"encoding/json"
	"github.com/tietang/assert"
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestGetMatchRouteTargetPath(t *testing.T) {
	Convey("基本匹配1", t, func() {
		r := &Router{}

		r.AddRoute(Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false})
		r.AddRoute(Route{Source: "/app1/v2/user", ServiceId: "app1", Target: "/app1/v2/user", StripPrefix: false})
		r.AddRoute(Route{Source: "/app2/**", ServiceId: "app2", Target: "/app20/**", StripPrefix: true})
		r.AddRoute(Route{Source: "/app3/**", ServiceId: "app3", StripPrefix: false})
		r.AddRoute(Route{Source: "/app4/**", ServiceId: "app4", StripPrefix: true})
		//
		Convey("精确匹配1", func() {
			route := r.GetMatchRoute(r.Routes[0].Source)
			So(r.Routes[0].Target, ShouldEqual, route.Target)
			p := r.GetMatchRouteTargetPath(r.Routes[0].Source)
			So(r.Routes[0].Target, ShouldEqual, p)

		})
		//
		Convey("精确匹配2", func() {
			route := r.GetMatchRoute(r.Routes[1].Source)
			So(r.Routes[1].Target, ShouldEqual, route.Target)
			p := r.GetMatchRouteTargetPath(r.Routes[1].Source)
			So(r.Routes[1].Target, ShouldEqual, p)
		})
		Convey("模糊匹配1", func() {
			sp := "/app2/v1/demo"
			actualTargetPath := r.GetMatchRouteTargetPath(sp)
			expPath := strings.TrimSuffix(r.Routes[2].Target, "/**") + strings.TrimPrefix(sp, "/app2")
			So(expPath, ShouldEqual, actualTargetPath)
		})
		Convey("模糊匹配2", func() {
			//
			sp := "/app3/v1/demo"
			actualTargetPath := r.GetMatchRouteTargetPath(sp)
			So(sp, ShouldEqual, actualTargetPath)
			//
		})
		Convey("模糊匹配3", func() {
			sp := "/app4/v1/demo"
			actualTargetPath := r.GetMatchRouteTargetPath(sp)
			expPath := strings.TrimPrefix(sp, "/app4")
			So(expPath, ShouldEqual, actualTargetPath)
			//
		})
	})
}

func TestGetMatchRoute(t *testing.T) {
	r := &Router{}

	r.AddRoute(Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false})
	r.AddRoute(Route{Source: "/app1/v2/user", ServiceId: "app1", Target: "/app1/v2/user", StripPrefix: false})
	r.AddRoute(Route{Source: "/app2/**", ServiceId: "app2", Target: "/app20/**", StripPrefix: true})
	r.AddRoute(Route{Source: "/app3/**", ServiceId: "app3", StripPrefix: false})
	r.AddRoute(Route{Source: "/app4/**", ServiceId: "app4", StripPrefix: true})
	//
	route := r.GetMatchRoute(r.Routes[0].Source)
	assert.Equal(t, r.Routes[0].ServiceId, route.ServiceId)
	assert.Equal(t, r.Routes[0].Target, route.Target)
	p := r.GetMatchRouteTargetPath(r.Routes[0].Source)
	assert.Equal(t, r.Routes[0].Target, p)

	//
	route = r.GetMatchRoute(r.Routes[1].Source)
	assert.Equal(t, r.Routes[1].ServiceId, route.ServiceId)
	assert.Equal(t, r.Routes[1].Target, route.Target)
	p = r.GetMatchRouteTargetPath(r.Routes[1].Source)
	assert.Equal(t, r.Routes[1].Target, p)

	sp := "/app2/v1/demo"
	route = r.GetMatchRoute(sp)
	assert.Equal(t, r.Routes[2].ServiceId, route.ServiceId)
	assert.Equal(t, r.Routes[2].Target, route.Target)
	actualTargetPath := r.GetMatchRouteTargetPath(sp)
	expPath := strings.TrimSuffix(r.Routes[2].Target, "/**") + strings.TrimPrefix(sp, "/app2")
	assert.Equal(t, expPath, actualTargetPath)
	//
	sp = "/app3/v1/demo"
	route = r.GetMatchRoute(sp)
	assert.Equal(t, r.Routes[3].ServiceId, route.ServiceId)
	assert.Equal(t, r.Routes[3].Target, route.Target)
	actualTargetPath = r.GetMatchRouteTargetPath(sp)
	assert.Equal(t, sp, actualTargetPath)
	//
	sp = "/app4/v1/demo"
	route = r.GetMatchRoute(sp)
	assert.Equal(t, r.Routes[4].ServiceId, route.ServiceId)
	assert.Equal(t, r.Routes[4].Target, route.Target)
	actualTargetPath = r.GetMatchRouteTargetPath(sp)
	expPath = strings.TrimPrefix(sp, "/app4")
	assert.Equal(t, expPath, actualTargetPath)
	//

}

func TestGetMatchRouteByJson(t *testing.T) {
	s := "[\n" +
		"    {\n" +
		"        \"ServiceId\": \"app1\",\n" +
		"        \"Source\": \"/app1/v1/user\",\n" +
		"        \"SourcePrefix\": \"\",\n" +
		"        \"SourceIsFuzzyMatch\": false,\n" +
		"        \"Target\": \"/v1/user\",\n" +
		"        \"TargetIsFuzzyMatch\": false,\n" +
		"        \"TargetPrefix\": \"\",\n" +
		"        \"StripPrefix\": false\n" +
		"    },\n" +
		"    {\n" +
		"        \"ServiceId\": \"app1\",\n" +
		"        \"Source\": \"/app1/info\",\n" +
		"        \"SourcePrefix\": \"\",\n" +
		"        \"SourceIsFuzzyMatch\": false,\n" +
		"        \"Target\": \"/info\",\n" +
		"        \"TargetIsFuzzyMatch\": false,\n" +
		"        \"TargetPrefix\": \"\",\n" +
		"        \"StripPrefix\": false\n" +
		"    },\n" +
		"    {\n" +
		"        \"ServiceId\": \"APP1\",\n" +
		"        \"Source\": \"/app1/**\",\n" +
		"        \"SourcePrefix\": \"/app1\",\n" +
		"        \"SourceIsFuzzyMatch\": true,\n" +
		"        \"Target\": \"/app1/**\",\n" +
		"        \"TargetIsFuzzyMatch\": true,\n" +
		"        \"TargetPrefix\": \"/app1\",\n" +
		"        \"StripPrefix\": false\n" +
		"    }" +
		"]"

	r := &Router{}
	routes := make([]Route, 3)
	json.Unmarshal([]byte(s), &routes)
	r.SetRoutes(routes)

	Convey("通过JSON路由模糊匹配", t, func() {
		sp := "/app1/v1/demo"
		route := r.GetMatchRoute(sp)
		So(route, ShouldNotBeNil)
		So(route.ServiceId, ShouldEqual, r.Routes[2].ServiceId)
		So(route.Target, ShouldEqual, r.Routes[2].Target)
		//assert.NotNil(t, route)
		//assert.Equal(t, r.Routes[2].ServiceId, route.ServiceId)
		//assert.Equal(t, r.Routes[2].Target, route.Target)
		actualTargetPath := r.GetMatchRouteTargetPath(sp)
		expPath := strings.TrimSuffix(r.Routes[2].Target, "/**") + strings.TrimPrefix(sp, "/app2")
		So(actualTargetPath, ShouldEqual, expPath)
	})

}
