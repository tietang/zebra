package zuul

import (
    "encoding/json"
    "strings"
    "testing"

    "github.com/tietang/assert"
    . "github.com/smartystreets/goconvey/convey"
)

func TestGetMatchRouteTargetPath(t *testing.T) {
    r := &Router{}

    r.AddRoute(Route{Source: "/app1/v1/user", App: "app1", Target: "/v1/user", StripPrefix: false})
    r.AddRoute(Route{Source: "/app1/v2/user", App: "app1", Target: "/app1/v2/user", StripPrefix: false})
    r.AddRoute(Route{Source: "/app2/**", App: "app2", Target: "/app20/**", StripPrefix: true})
    r.AddRoute(Route{Source: "/app3/**", App: "app3", StripPrefix: false})
    r.AddRoute(Route{Source: "/app4/**", App: "app4", StripPrefix: true})
    //
    route := r.GetMatchRoute(r.Routes[0].Source)
    assert.Equal(t, r.Routes[0].Target, route.Target)
    p := r.GetMatchRouteTargetPath(r.Routes[0].Source)
    assert.Equal(t, r.Routes[0].Target, p)
    //
    route = r.GetMatchRoute(r.Routes[1].Source)
    assert.Equal(t, r.Routes[1].Target, route.Target)
    p = r.GetMatchRouteTargetPath(r.Routes[1].Source)
    assert.Equal(t, r.Routes[1].Target, p)

    sp := "/app2/v1/demo"
    actualTargetPath := r.GetMatchRouteTargetPath(sp)
    expPath := strings.TrimSuffix(r.Routes[2].Target, "/**") + strings.TrimPrefix(sp, "/app2")
    assert.Equal(t, expPath, actualTargetPath)
    //
    sp = "/app3/v1/demo"
    actualTargetPath = r.GetMatchRouteTargetPath(sp)
    assert.Equal(t, sp, actualTargetPath)
    //
    sp = "/app4/v1/demo"
    actualTargetPath = r.GetMatchRouteTargetPath(sp)
    expPath = strings.TrimPrefix(sp, "/app4")
    assert.Equal(t, expPath, actualTargetPath)
    //

}

func TestGetMatchRoute(t *testing.T) {
    r := &Router{}

    r.AddRoute(Route{Source: "/app1/v1/user", App: "app1", Target: "/v1/user", StripPrefix: false})
    r.AddRoute(Route{Source: "/app1/v2/user", App: "app1", Target: "/app1/v2/user", StripPrefix: false})
    r.AddRoute(Route{Source: "/app2/**", App: "app2", Target: "/app20/**", StripPrefix: true})
    r.AddRoute(Route{Source: "/app3/**", App: "app3", StripPrefix: false})
    r.AddRoute(Route{Source: "/app4/**", App: "app4", StripPrefix: true})
    //
    route := r.GetMatchRoute(r.Routes[0].Source)
    assert.Equal(t, r.Routes[0].App, route.App)
    assert.Equal(t, r.Routes[0].Target, route.Target)
    p := r.GetMatchRouteTargetPath(r.Routes[0].Source)
    assert.Equal(t, r.Routes[0].Target, p)

    //
    route = r.GetMatchRoute(r.Routes[1].Source)
    assert.Equal(t, r.Routes[1].App, route.App)
    assert.Equal(t, r.Routes[1].Target, route.Target)
    p = r.GetMatchRouteTargetPath(r.Routes[1].Source)
    assert.Equal(t, r.Routes[1].Target, p)

    sp := "/app2/v1/demo"
    route = r.GetMatchRoute(sp)
    assert.Equal(t, r.Routes[2].App, route.App)
    assert.Equal(t, r.Routes[2].Target, route.Target)
    actualTargetPath := r.GetMatchRouteTargetPath(sp)
    expPath := strings.TrimSuffix(r.Routes[2].Target, "/**") + strings.TrimPrefix(sp, "/app2")
    assert.Equal(t, expPath, actualTargetPath)
    //
    sp = "/app3/v1/demo"
    route = r.GetMatchRoute(sp)
    assert.Equal(t, r.Routes[3].App, route.App)
    assert.Equal(t, r.Routes[3].Target, route.Target)
    actualTargetPath = r.GetMatchRouteTargetPath(sp)
    assert.Equal(t, sp, actualTargetPath)
    //
    sp = "/app4/v1/demo"
    route = r.GetMatchRoute(sp)
    assert.Equal(t, r.Routes[4].App, route.App)
    assert.Equal(t, r.Routes[4].Target, route.Target)
    actualTargetPath = r.GetMatchRouteTargetPath(sp)
    expPath = strings.TrimPrefix(sp, "/app4")
    assert.Equal(t, expPath, actualTargetPath)
    //

}

func TestGetMatchRouteByJson(t *testing.T) {
    s := "[\n" +
            "    {\n" +
            "        \"App\": \"app1\",\n" +
            "        \"Source\": \"/app1/v1/user\",\n" +
            "        \"SourcePrefix\": \"\",\n" +
            "        \"SourceIsFuzzyMatch\": false,\n" +
            "        \"Target\": \"/v1/user\",\n" +
            "        \"TargetIsFuzzyMatch\": false,\n" +
            "        \"TargetPrefix\": \"\",\n" +
            "        \"StripPrefix\": false\n" +
            "    },\n" +
            "    {\n" +
            "        \"App\": \"app1\",\n" +
            "        \"Source\": \"/app1/info\",\n" +
            "        \"SourcePrefix\": \"\",\n" +
            "        \"SourceIsFuzzyMatch\": false,\n" +
            "        \"Target\": \"/info\",\n" +
            "        \"TargetIsFuzzyMatch\": false,\n" +
            "        \"TargetPrefix\": \"\",\n" +
            "        \"StripPrefix\": false\n" +
            "    },\n" +
            "    {\n" +
            "        \"App\": \"APP1\",\n" +
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
        So(route.App, ShouldEqual, r.Routes[2].App)
        So(route.Target, ShouldEqual, r.Routes[2].Target)
        //assert.NotNil(t, route)
        //assert.Equal(t, r.Routes[2].App, route.App)
        //assert.Equal(t, r.Routes[2].Target, route.Target)
        actualTargetPath := r.GetMatchRouteTargetPath(sp)
        expPath := strings.TrimSuffix(r.Routes[2].Target, "/**") + strings.TrimPrefix(sp, "/app2")
        So(actualTargetPath, ShouldEqual, expPath)
    })


    //

}
