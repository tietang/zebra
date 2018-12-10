package main

import (
    "github.com/tietang/zebra"
    //"github.com/tietang/logrus-prefixed-formatter"
    _ "net/http/pprof"
)

func main() {

    //server := doorsill.New("proxy.properties")
    address := "127.0.0.1:8500"
    root := "zebra/bootstrap"
    server := zebra.NewByConsulKeyValue(address, root)
    //server := New("/Users/tietang/Documents/git.oschina/r_app/src/github.com/tietang/zebra/example/proxy.properties")
    server.Start()
}
