package router

import (
	"fmt"
	"github.com/samuel/go-zookeeper/zk"
	"github.com/tietang/props/kvs"
	"github.com/tietang/props/zk"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"testing"
)

var kv map[string]string

func init() {

	kv = make(map[string]string)
}
func TestNewZkRouteSource(t *testing.T) {
	dir, _ := os.Getwd()
	//cmd := filepath.Join(dir, "mock.jar")
	//var c *exec.Cmd
	////go func() {
	//c = exec.Command(cmd, "run")
	//fmt.Println(c.Args)
	//stdout, err := c.StdoutPipe()
	//
	//if err != nil {
	//	fmt.Println(err)
	//}
	//
	//c.Start()
	//reader := bufio.NewReader(stdout)
	//
	////实时循环读取输出流中的一行内容
	//for {
	//	line, err2 := reader.ReadString('\n')
	//	if err2 != nil || io.EOF == err2 {
	//		break
	//	}
	//	if strings.contains(line, "started") {
	//		break
	//
	//	}
	//	fmt.Println(line)
	//}

	//}()
	file := filepath.Join(dir, "zk_test.properties")
	ps := kvs.NewPropertiesConfigSource(file)
	p := NewZookeeperRouteSource(ps)
	p.Init()
	size := 3

	p.Init()
	p.Start()
	initData(p.zkConn, p.zkContexts, size)
	Convey("1", t, func() {

		for _, route := range p.routes {
			fmt.Println(route)
		}
		So(p.routes, ShouldNotBeNil)
		So(len(p.routes), ShouldEqual, size)

	})

	//c = exec.Command(cmd, "stop")
	//c.Wait()
}

func initData(zkConn *zk.Conn, contexts []string, size int) {
	for i := 0; i < size; i++ {
		key := "key-" + strconv.Itoa(i)
		strSeq := strconv.Itoa(i)
		value := "/app" + strSeq + "/**=app" + strSeq + ""
		p := path.Join(contexts[0], key) //  + "/" + KEY_ROUTES_NODE_PREFIX + "/" + key
		vkey := "app.xx." + key
		if !kvs.ZkExits(zkConn, p) {
			_, err := kvs.ZkCreateString(zkConn, p, value)
			fmt.Println(key, ": ", value)
			if err == nil {
				//log.Println(path)
				kv[vkey] = value
			}
			//log.Println(err)
		} else {
			kv[vkey] = value
		}

	}

	//fmt.Println(len(kv))
}
