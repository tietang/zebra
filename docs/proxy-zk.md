

### 通过zookeeper配置来启动proxy server



函数：NewHttpProxyServerByZookeeper(zkUrls []string, contexts []string, rootPath string) *HttpProxyServer

zkUrls zookeeper连接字符串

contexts zookeeper 路径上下文相对路径

rootPath zookeeper 绝对root路径

### 配置

#### server.port
默认是8002，配置proxy server端口
#### eureka.server.urls
eureka访问地址，多个用","分隔,比如：
http://192.168.1.11:8761/eureka,http://192.168.1.12:8761/eureka,http://192.168.1.13:8761/eureka

#### routes

##### 手动路由配置，在routes节点下配置key/value：

- key 为任意，能说明区分所配置的项目
- value，标准的java properties文件形式的内容，
	- 其格式为：源路径=服务名称,[目标路径],[true|false是否截取前缀]
	- 源路径：必须，可以`前缀/**`格式来模糊路由;可以是具体明确的路径
	- 服务名称：必须，如果配置服务发现，其名称要和服务发现中的名称一致；否则，为任一能明确区分唯一的描述
	- 目标路径：可选，规则同源路径
	- 是否截取前缀：表示是否截取模糊匹配形式下的`前缀`截取

例子(用yaml格式来描述)：

```yaml
routes:
  api: >
    app1:,/app,/app,false
  user: >
    /user/**,user,/user/**,false
    /user1/**,user,/user1/**,false
    /user2/**,user,/user/**
    /user3/**,user
    /user1/v1/users,user
    /user2/v1/users,user,/user/v1/users
    /user3/v1/users,user,/user/v1/users,true	
```

#### example:

```go
urls := []string{"172.16.1.248:2181"}
contexts := []string{"apps"}
proxy := zuul.NewHttpProxyServerByZookeeper(urls, contexts,"/configs/")
proxy.DefaultRun()
    
```

```

/configs/apps/eureka/server=urls=http://127.0.0.1:8761/eureka
/configs/apps/hystrix/default=ErrorPercentThreshold=50
/configs/apps/hystrix/default=MaxConcurrentRequests=100
/configs/apps/hystrix/default=RequestVolumeThreshold=20
/configs/apps/hystrix/default=SleepWindow=5000
/configs/apps/hystrix/default=Timeout=6000
/configs/apps/routes=router=##app-order\n/order/v1/create/=apporder,/v1/order/create
/configs/apps/server=port=19002
```


## Hystrix 配置
格式：hystrix.`hystrix command key`.`参数`
对于服务来说：`hystrix command key`是服务名称的小写

### 默认配置：
如果某个服务没有配置，则默认配置会起作用

```
hystrix.default.Timeout=6000
hystrix.default.MaxConcurrentRequests=100
hystrix.default.RequestVolumeThreshold=20
hystrix.default.SleepWindow=5000
hystrix.default.ErrorPercentThreshold=50

```
