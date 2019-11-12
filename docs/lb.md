# 负载均衡配置

```ini

[lb]
default.balancer.name = WeightRobinRound
default.max.fails = 3
default.fail.time.window = 10s
default.fail.sleep.mode = seq
default.fail.sleep.x = 1,1,2,3,5,8,13,21
default.fail.sleep.max = 60s
log.selected.enabled = true
log.selected.all = false

```

`default`为全局默认标识，对于具体的服务配置可将`default`替换为服务id即可单独来配置该服务的负载均衡器，参考[routes 文件配置规则](<routes_config.md>).


### balancer.name 

全局默认的负载均衡器算法名称，目前支持以下负载均衡算法：

- WeightRobinRound 加权轮询
- random 随机选取
- round 简单轮询
- hash 一致性hash
- fibonacci 基于平均响应时间斐波纳契算法

 默认为WeightRobinRound.

#### max.fails 

最大失败次数，在失败时间窗口`fail.time.window`给定服务的某个节点内达到`max.fails`设定的值，则负载均衡在`fail.time.window×fail.sleep.x`时间内不再选中该节点，在下一时间窗口再次重试。

#### fail.time.window 

失败检测时间窗口

#### fail.sleep.mode

失败时间窗口模式：fixed（固定倍数）和（seq）序列

- fixed：fail.sleep.x的值应该是1个值
- seq：fail.sleep.x为逗号分隔的数字，代表从第一个窗口开始，后面每次sleep的倍数，默认配置为一个斐波纳契数列：1,1,2,3,5,8,13,21

#### fail.sleep.x 

失败后睡眠时间窗口倍数，根据fail.sleep.mode来配置，参考fail.sleep.mode中的说明。
该配置为一个数字序列或单值。

#### fail.sleep.max

失败熔断窗口最大时间，-1 或者 小于`default.fail.time.window`不对其起效；如果失败熔断窗口时间大于该值，则后面失败时间窗口模式被设定为fixed，后续的`fail.sleep.x`倍数序列不在起效。


#### log.selected.enabled = true

是否启用负载均衡选择实例的日志输出。

#### log.selected.all = false

是否启用负载均衡选择实例的日志输出完整模式。

- true 输出整个instacne信息
- false 只输出实例id和address

