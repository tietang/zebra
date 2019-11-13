package httpclient

import (
	"github.com/tietang/props/kvs"
	"net/http"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDo(t *testing.T) {
	psc := kvs.NewPropertiesConfigSource("httpclient.properties")

	hcs := NewHttpClients(psc)
	name := "test"
	url := "http://www.baidu.com/"
	Convey("get", t, func() {

		req, _ := http.NewRequest("GET", url, nil)
		res, err := hcs.Do(name, req)
		So(res, ShouldNotBeNil)
		So(res.Body, ShouldNotBeNil)
		So(err, ShouldBeNil)
	})

}
