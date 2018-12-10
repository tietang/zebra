#K8S 动态服务发现和动态路由

##规范：

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
