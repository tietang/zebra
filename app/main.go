package main

import (
	"gitee.com/tietang/terrace-go/base/ilogrus"
	"gitee.com/tietang/terrace-go/boot"
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "github.com/tietang/zebra"
	"os"
	"path/filepath"
)

func main() {
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("proxy.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	ilogrus.InitLog(conf)
	app := boot.New(conf)
	app.RapidStart()
}

func GetCurrentFilePath(fileName string, skip int) string {
	//获取当前函数调用对应的文件
	dir, _ := os.Getwd()
	file := filepath.Join(dir, fileName)
	return file
}
