# 路由规则配置

路由规则支持4种方式的配置：

- ini格式的配置文件
- zookeeper 作为存储源
- consul 作为存储源
- mysql 作为存储源



## routes 文件配置规则

配置规则适用于除mysql方式的其他所有存储源。

### 格式

```ini

#服务配置
[{serviceId}]

lb.name = WeightRobinRound
max.fails = 3
fail.time.window = 10s
fail.sleep.mode = seq|fixed
fail.sleep.x = 1,1,2,3,5,8,13,21|1
routes.strip.prefix = false


#服务路由规则配置
[{serviceId}.routes]

`source path`=`target path`,`isStripPrefix`

```




#### [{serviceId}]
serviceId为服务id或者服务名称，对应于服务发现中的服务ID：

- Eureka中是appName, `application.name`或者`instance.app`
- consul中的ServiceID或者ServiceName
- kubernetes中的metadata.name


#### 负载均衡配置

```
lb.name = WeightRobinRound
max.fails = 3
fail.time.window = 10s
fail.sleep.mode = seq|fixed
fail.sleep.x = 1,1,2,3,5,8,13,21|1
```

详细配置参考[负载均衡配置](<lb.md>).



#### routes.strip.prefix  

当路由规则为模糊匹配规则时，是否截去模糊匹配前缀。


#### 路由配置

[{serviceId}.routes]节点下配置路由规则：`source path`=`target path`,`isStripPrefix`

##### source path

源匹配路径，支持精确匹配和模糊匹配2种规则：

- 模糊匹配，在路径最后面用`**`来，例如: `/api/users/**`
- 精确匹配，完整的除域名之外的path，，例如: `/api/users/name/tietang`

##### target path

目标转发路径, 支持精确转发和模糊拼接转发2种规则：

- 精确，匹配到源路径，直接替换为目标路径
- 模糊拼接，匹配到源路径，
	- 源路径为模糊匹配
		- isStripPrefix=true, 替换目标路径中的`**`为将源路径中的`**`
		- isStripPrefix=false, 替换目标路径中的`**`为将源路径完整的path
	- 源路径为精确匹配，替换目标路径中的`**`为将源路径完整的path
		
##### isStripPrefix

是否截去前缀，只有`source path`为模糊匹配时。




  
例如：

```ini
[user]
max.fails = 3
fail.time.window = 10s
fail.sleep.x = 1

[user.routes]

/user/v1/**=/v1/**
/user/v2/users=/v2/users,false

[order]

[order.routes]
/order/v1/**=/order/v1/**
/order/v2/list=/order/v2/list,false


```




### ini格式的配置文件

在zebra配置文件(默认为proxy.ini)中找到`[ini]` 节点，如下：

```ini

[ini]
#ini 文件的静态路由规则配置
routes.enabled = true
#routes.dir=/Users/tietang/my/gitcode/r_app/src/github.com/tietang/zebra/example/routes

```

#### routes.enabled 

启用ini格式文件的routes配置，默认为false。

#### routes.dir 

routes规则配置文件目录，默认为`routes`。
如果不配置或者注释，则使用默认值，默认是读取程序运行同目录下的`routes`文件夹，会读取该文件夹下的所有文件，文件格式参考"routes 文件配置规则"。
 
如果配置目录是绝对路径，则会直接读取路径文件夹。
如果配置目录是相对路径，则会在程序运行同目录下查找所配置的目录，并读取。

`routes.dir`中文件名称无实际含义，可以根据自己需求来定义，内容格式参考"routes 文件配置规则"，可以在一个文件中配置多个serviceId，也可以分为多个文件配置。


### zookeeper 作为存储源


在zebra配置文件(默认为proxy.ini)中找到`[zk]` 节点，如下：

```ini

[zk]
## zookeeper 静态路由规则配置
enabled = true
##zk连接字符串，多个逗号分隔:192.168.1.2:2181,192.168.1.3:2181,192.168.1.4:2181
;conn.urls = 127.0.0.1:2181
conn.urls = 172.16.1.248:2181
timeout = 3s
#zk跟节点
root = /zebra

#多个逗号分割
routes.contexts = routes
#notify-node,all
routes.watch.type = all

```

#### enabled 

是否启用zookeeper作为配置源。

#### conn.urls 

zookeeper连接字符串，多个逗号分隔，例如: `192.168.1.2:2181,192.168.1.3:2181,192.168.1.4:2181`。

#### timeout

zookeeper连接超时时间，默认为3s。

#### root

zookeeper节点的root节点，默认为：`/zebra`

#### routes.contexts 

routes 配置的上下文路径，该路径为相对路径，配合root一起使用，例如：

```
root=/zebra
routes.contexts= routes
那么routes.contexts的实际绝对路径是：/zebra/routes

```

#### routes.watch.type 

routes的节点监听类型，支持2种：notify-node,all

##### notify-node

会在`routes.contexts`下创建一个名称为`notify`的节点，节点值为任意，当修改routes规则节点后，只需要修改`notify`节点的值和之前的值不一样，触发值变事件，zebra中会自动重新读取routes，并启用。

##### all

直接监听`routes.contexts`下所有的子节点，只要子节点有变化，触发值变事件，zebra中会自动重新读取变化的节点并更新节点内容，并启用。


routes节点配置key任意命名，value参考“routes 文件配置规则”，下面是配置图示：

![](<imgs/zk-routes-1.png>)
![](<imgs/zk-routes.png>)

### consul 作为存储源



consul作为配置源，其配置方式和zookeeper类似。另外consul支持服务发现，参考服务发现。

```ini
[consul]
enabled = false
address = 172.16.1.248:8500
# routes
routes.enabled = true
routes.root = zebra/routes
```




![](<imgs/consul-routes.png>)



### mysql 作为存储源

