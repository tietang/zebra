package httpclient

import (
    "github.com/tietang/props/kvs"
    "net/http"
    "testing"
)

func TestDo(t *testing.T) {
    psc := kvs.NewPropertiesConfigSource("/Users/tietang/my/gitcode/r_app/src/github.com/tietang/go-zuul/httpclient/httpclient.properties")

    hcs := NewHttpClients(psc)
    name := "test"
    url := "http://www.baidu.com/"
    Convey("get", t, func() {

        req, _ := http.NewRequest("GET", url, nil)
        res, body, err := hcs.Do(name, req)
        So(res, ShouldNotBeNil)
        So(body, ShouldNotBeNil)
        So(err, ShouldBeNil)
    })

}
