package main

import (
    "github.com/tietang/zebra/tools"
)

func main() {
    file := "/Users/tietang/my/gitcode/r_app/src/github.com/tietang/zebra/example/proxy.ini"
    address := "172.16.1.248:8500"
    root := "zebra/bootstrap"
    //tools.FileToConsulKeyValue(file, address, root)
    tools.IniFileToConsulProperties(file, address, root)
    //tools.FileToZookeeperKeyValue(file, "127.0.0.1:2181", "/"+root)
    tools.IniFileToZookeeperProperties(file, "172.16.1.248:2181", "/"+root)

}
