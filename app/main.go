package main

import (
	"github.com/tietang/props/ini"
	"github.com/tietang/props/kvs"
	_ "github.com/tietang/zebra"
	"github.com/tietang/zebra/infra"
	"github.com/tietang/zebra/infra/base"
	"os"
	"path/filepath"
)

func main() {
	//获取程序运行文件所在的路径
	file := kvs.GetCurrentFilePath("proxy.ini", 1)
	//加载和解析配置文件
	conf := ini.NewIniFileCompositeConfigSource(file)
	base.InitLog(conf)
	app := infra.New(conf)
	app.Start()
}

func GetCurrentFilePath(fileName string, skip int) string {
	//获取当前函数Caller reports，取得当前调用对应的文件
	dir, _ := os.Getwd()
	file := filepath.Join(dir, fileName)
	return file
}
