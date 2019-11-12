package router

import (
    "testing"
)

func TestAll(t *testing.T) {
    Convey("初始化", t, func() {

        Convey("One to One 1", func() {
            r := Route{Source: "/app1/v1/user", ServiceId: "app1", Target: "/v1/user", StripPrefix: false}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeFalse)
            So(r.SourcePrefix, ShouldBeBlank)
            So(r.TargetIsFuzzyMatch, ShouldBeFalse)
            So(r.TargetPrefix, ShouldBeBlank)
            So(r.StripPrefix, ShouldBeFalse)
            targetPath := r.GetRouteTargetPath(r.Source)
            So(targetPath, ShouldEqual, r.Target)

        })

        Convey("One to One 2", func() {
            r := Route{Source: "/app1/v2/user", ServiceId: "app1", Target: "/app1/v2/user", StripPrefix: true}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeFalse)
            So(r.SourcePrefix, ShouldBeBlank)
            So(r.TargetIsFuzzyMatch, ShouldBeFalse)
            So(r.TargetPrefix, ShouldBeBlank)
            So(r.StripPrefix, ShouldBeTrue)
            targetPath := r.GetRouteTargetPath(r.Source)
            So(targetPath, ShouldEqual, r.Target)
        })
        Convey("Fuzzy source&target 1", func() {
            r := Route{Source: "/app2/**", ServiceId: "app2", Target: "/app20/**", StripPrefix: true}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeTrue)
            So(r.SourcePrefix, ShouldEqual, "/app2")
            So(r.TargetIsFuzzyMatch, ShouldBeTrue)
            So(r.TargetPrefix, ShouldEqual, "/app20")
            So(r.StripPrefix, ShouldBeTrue)
            targetPath := r.GetRouteTargetPath("/app2/v1/users")
            So(targetPath, ShouldEqual, "/app20/v1/users")
        })

        Convey("Fuzzy source&target 2", func() {
            r := Route{Source: "/app2/**", ServiceId: "app2", Target: "/app20/**", StripPrefix: false}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeTrue)
            So(r.SourcePrefix, ShouldEqual, "/app2")
            So(r.TargetIsFuzzyMatch, ShouldBeTrue)
            So(r.TargetPrefix, ShouldEqual, "/app20")
            So(r.StripPrefix, ShouldBeFalse)
            targetPath := r.GetRouteTargetPath("/app2/v1/users")
            So(targetPath, ShouldEqual, "/app20/app2/v1/users")
        })

        Convey("Fuzzy source&target 3", func() {
            r := Route{Source: "/app3/**", ServiceId: "app3", StripPrefix: false}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeTrue)
            So(r.SourcePrefix, ShouldEqual, "/app3")
            So(r.TargetIsFuzzyMatch, ShouldBeFalse)
            So(r.TargetPrefix, ShouldBeBlank)
            So(r.StripPrefix, ShouldBeFalse)
            targetPath := r.GetRouteTargetPath("/app3/v1/users")
            So(targetPath, ShouldEqual, "/app3/v1/users")
        })
        Convey("Fuzzy source&target 4", func() {
            r := Route{Source: "/app4/**", ServiceId: "app4", StripPrefix: true}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeTrue)
            So(r.SourcePrefix, ShouldEqual, "/app4")
            So(r.TargetIsFuzzyMatch, ShouldBeFalse)
            So(r.TargetPrefix, ShouldBeBlank)
            So(r.StripPrefix, ShouldBeTrue)
            targetPath := r.GetRouteTargetPath("/app4/v1/users")
            So(targetPath, ShouldEqual, "/v1/users")

        })

        //registeredPath=/portal_sync_io/{sub_path:path} route=
        // {ServiceId:PORTAL_SYNC_IO Source:/portal_sync_io/** SourcePrefix: SourceIsFuzzyMatch:false Target: TargetIsFuzzyMatch:false TargetPrefix: StripPrefix:true}
        Convey("Fuzzy source&target example", func() {
            r := Route{ServiceId: "PORTAL_SYNC_IO", Source: "/portal_sync_io/**", StripPrefix: true}
            r.Init()
            So(r.SourceIsFuzzyMatch, ShouldBeTrue)
            So(r.SourcePrefix, ShouldEqual, "/portal_sync_io")
            So(r.TargetIsFuzzyMatch, ShouldBeFalse)
            So(r.TargetPrefix, ShouldBeBlank)
            So(r.StripPrefix, ShouldBeTrue)
            targetPath := r.GetRouteTargetPath("/portal_sync_io/v1/users")
            So(targetPath, ShouldEqual, "/v1/users")

        })
    })
}
