package tools

import (
    "errors"
    "fmt"
    "github.com/hashicorp/consul/api"
    "github.com/samuel/go-zookeeper/zk"
    "github.com/tietang/go-utils"
    "github.com/tietang/props/kvs"
    "github.com/tietang/props/zk"
    "gopkg.in/ini.v1"
    "os"
    "path"
    "path/filepath"
    "strings"
    "time"
)

func FileToConsulKeyValue(file string, address, root string) {
    if !filepath.IsAbs(file) {
        dir, err := os.Getwd()
        utils.Panic(err)
        file = filepath.Join(dir, file)
    }
    isExists, err := utils.PathExists(file)
    utils.Panic(err)
    if !isExists {
        panic(errors.New("file is not exists: " + file))
    }
    ext := filepath.Ext(file)
    var conf kvs.ConfigSource
    if strings.Contains(ext, "prop") {
        conf = kvs.NewPropertiesConfigSource(file)
    } else {
        conf = kvs.NewIniFileConfigSource(file)
    }

    config := api.DefaultConfig()
    config.Address = address
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }
    kv := client.KV()
    wq := &api.WriteOptions{}
    keys := conf.Keys()
    for _, key := range keys {

        keyFull := path.Join(root, strings.Replace(key, ".", "/", -1))
        value := conf.GetDefault(key, "")
        kvp := &api.KVPair{
            Key:   keyFull,
            Value: []byte(value),
        }
        kv.Put(kvp, wq)
        fmt.Println(kvp.Key, "=", string(kvp.Value))
    }
}

func IniFileToConsulProperties(file string, address, root string) {
    if !filepath.IsAbs(file) {
        dir, err := os.Getwd()
        utils.Panic(err)
        file = filepath.Join(dir, file)
    }
    isExists, err := utils.PathExists(file)
    utils.Panic(err)
    if !isExists {
        panic(errors.New("file is not exists: " + file))
    }
    iniFile, err := ini.Load(file)

    config := api.DefaultConfig()
    config.Address = address
    client, err := api.NewClient(config)
    if err != nil {
        panic(err)
    }
    kv := client.KV()
    wq := &api.WriteOptions{}
    sections := iniFile.Sections()
    for _, section := range sections {
        name := section.Name()
        keyFull := path.Join(root, strings.Replace(name, ".", "/", -1))
        values := make([]string, 0)
        for _, key := range section.Keys() {
            line := key.Name() + " = " + key.Value()
            values = append(values, line)
        }
        fmt.Println()
        fmt.Println(name, ": ")
        value := strings.Join(values, "\n")
        kvp := &api.KVPair{
            Key:   keyFull,
            Value: []byte(value),
        }
        kv.Put(kvp, wq)
        fmt.Println(value)
    }

}

func FileToZookeeperKeyValue(file string, urls, root string) {
    if !filepath.IsAbs(file) {
        dir, err := os.Getwd()
        utils.Panic(err)
        file = filepath.Join(dir, file)
    }
    isExists, err := utils.PathExists(file)
    utils.Panic(err)
    if !isExists {
        panic(errors.New("file is not exists: " + file))
    }
    ext := filepath.Ext(file)
    var conf kvs.ConfigSource
    if strings.Contains(ext, "prop") {
        conf = kvs.NewPropertiesConfigSource(file)
    } else {
        conf = kvs.NewIniFileConfigSource(file)
    }

    conn, ch, err := zk.Connect([]string{urls}, 2*time.Second)
    if err != nil {
        panic(err)
    }
    for {
        event := <-ch
        fmt.Println(event)
        if event.State == zk.StateConnected {
            break
        }
    }

    keys := conf.Keys()
    for _, key := range keys {

        keyPath := path.Join(root, strings.Replace(key, ".", "/", -1))
        value := conf.GetDefault(key, "")

        if !kvs.ZkExits(conn, keyPath) {
            _, err := kvs.ZkCreateString(conn, keyPath, value)
            if err == nil {
                //log.Println(path)
            }
            //fmt.Println(v)
        }
        fmt.Println(keyPath, "=", value)
    }
}

func IniFileToZookeeperProperties(file, urls, root string) {
    if !filepath.IsAbs(file) {
        dir, err := os.Getwd()
        utils.Panic(err)
        file = filepath.Join(dir, file)
    }
    isExists, err := utils.PathExists(file)
    utils.Panic(err)
    if !isExists {
        panic(errors.New("file is not exists: " + file))
    }
    iniFile, err := ini.Load(file)

    conn, ch, err := zk.Connect([]string{urls}, 2*time.Second)
    if err != nil {
        panic(err)
    }
    for {
        event := <-ch
        fmt.Println(event)
        if event.State == zk.StateConnected {
            break
        }
    }
    sections := iniFile.Sections()
    for _, section := range sections {
        name := section.Name()
        keyPath := path.Join(root, strings.Replace(name, ".", "/", -1))
        values := make([]string, 0)
        for _, key := range section.Keys() {
            line := key.Name() + " = " + key.Value()
            values = append(values, line)
        }
        fmt.Println()
        fmt.Println(name, ": ")
        val := strings.Join(values, "\n")
        fmt.Println(val)
        if !kvs.ZkExits(conn, keyPath) {
            _, err := kvs.ZkCreateString(conn, keyPath, val)
            if err == nil {
                //log.Println(path)
            }
            //fmt.Println(v)
        }

    }

}
