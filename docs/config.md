# 配置

## 启动配置

`bootstrap.`开头的属性用于配置源的配置，支持`properties文件或硬编码`、`zookeeper key/value`、`consul key/valu` 3中配置源。
仅用于通过配置文件启动server的方式。


### zookeeper作为配置源

- bootstrap.zk.enabled 值为true,false

- bootstrap.zk.urls=127.0.0.1:2181

- bootstrap.zk.timeout=10s

- bootstrap.zk.root=/zebra/bootstrap
