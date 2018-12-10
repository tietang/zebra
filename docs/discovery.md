# 动态服务发现

目前只支持以下3种服务发现：

- Eureka
- Consul
- Kubernetes

3种服务发现可以同时支持，启用服务发现后，zebra会自动从发现服务器拉取注册表信息，默认会把所有的服务注册到路由规则和负载均衡中，默认将serviceId转换为小写作为url前缀来路由，例如：
User -> /user/**，所有/user/请求都会转发到User服务，并负载均衡到注册发现的服务实例列表。

在zebra中，目前支持的3种服务发现都是是基于客户端服务发现思路构建的，在zebra中会有一个服务发现的注册表副本，服务发现的可用性不会直接影响到Zebra，服务发现的变更会有一定的延迟，延迟和zebra中配置的刷新间隔和服务发现本身的机制有关。


## Eureka服务发现

关于Eureka参考：

- [Eureka 官网](<https://github.com/netflix/eureka>)
- [Spring Cloud Netflix](<http://cloud.spring.io/spring-cloud-static/Dalston.SR3/#_spring_cloud_netflix>)

启用eureka后，zebra会自动从Eureka Server拉取注册表信息，默认会把所有的服务注册到路由规则和负载均衡中；默认将appName转换为小写作为url前缀来路由，例如：User-> /user/**。


### 配置

在zebra配置文件(默认为proxy.ini)中找到`[eureka]` 节点，如下：

```
[eureka]
# eureka 服务发现：动态路由和负载均衡
server.enabled = false
## 多个集群用逗号分隔
;server.cluster = cluster1,cluster2
#server.urls=http://eureka.didispace.com/eureka
server.urls = http://172.16.1.248:8761/eureka
discovery.interval = 10s

```

#### server.enabled

是否启用Eureka 服务发现

#### server.urls

eureka Server的url列表，多个用逗号分割，目前不支持多集群。

#### discovery.interval

eureka注册表更新间隔，默认10s


## Consul服务发现

### 配置

在zebra配置文件(默认为proxy.ini)中找到`[consul]` 节点，如下：

```
[consul]
enabled = false
address = 172.16.1.248:8500
## consul 服务发现：动态路由和负载均衡
discovery.enabled = true
#discovery.address=127.0.0.1:8500
discovery.address = ${consul.address}
discovery.interval = 10s

```

#### enabled

是否启用Consul支持。

#### address

Consul服务地址，格式为：IP:PORT

#### discovery.enabled 

是否启用服务发现。

#### discovery.address 

consul服务发现地址，默认是consul支持的address：{consul.address}

#### discovery.interval

Consul服务发现注册表更新间隔，默认10s


## K8S 动态服务发现和动态路由

### 规范：

在k8s控制台创建lables标签，key为route_prefix，值为url pattern的前缀，由于k8s label值只支持字母、数字和`.-_`字符，用`.`来替换URL中的`/`表示。

格式为：

```[pattern1].[pattern2].[patternn][.]```

相当于

```/[pattern1]/[pattern2]/[patternn][/**]```


如果包含最后一个`.`则表示为模糊匹配，例如:

- `pattern1.`         相当于`/pattern1/**`

- `pattern1.pattern2` 相当于`/pattern1/pattern2`

- `pattern1.pattern2.`相当于`/pattern1/pattern2/**`


例如：

应用名称：user
```properties
route_prefix=api.user.
//route_prefix=/api/user/**
```
表示请求URL path中`/api/user/`开头的API都会转发到 user服务。


### 配置

在zebra配置文件(默认为proxy.ini)中找到`[k8s]` 节点，如下：

```

[k8s]
# Kubernetes服务发现支持
enabled = false
urls = http://172.16.30.112:8080
discovery.interval = 10s

```

#### enabled 

是否启用服务发现。

#### urls 

k8s服务发现地址

#### discovery.interval

服务发现注册表更新间隔，默认10s

