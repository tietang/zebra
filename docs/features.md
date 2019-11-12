# 动态网关



一个基于golang动态网关，提供动态路由和负载均衡，容错。

##### 特征：

- 支持3种启动配置
	- 通过ini文件
	- 通过zookeeper
	- 通过consul
- 支持多种配置管理
	- ini文件
	- zookeeper
	- consul
	- sql
- 支持多种服务发现
	- Eureka
	- Consul
	- Kubernetes
	- zookeeper，规划中...
	- etcd，规划中...
- 基于服务发现的动态路由
- 负载均衡
	- 简单轮询
	- 加权轮询
	- 随机
	- 一致性hash
	- 基于响应时间使用fibonacci加权轮询
- 简单监控
- 隔离降级&限流
	- hystrix熔断
	- 失败次数算法
- metrics
