## anser是什么 ##

* 一个命令行工具，用于辅助基于微服务的server开发
* 一个rpc框架，集成grpc以及其他相关辅助服务
* 一组辅助组件，用于提供基础服务(http-gateway, grpc-gateway, anser-config，anser-monitor等)

### anser client ###

支持的命令有：
* setup: 安装开发用的k8s以及anser使用的辅助容器
* init-service: 生成基础的服务框架

### anser服务框架 ###

* 基于gRPC
* 集成gops
* 服务发现
* promethes监控
* jaeger跟踪？
* 分布式log(ELK?)
* envoy（service mesh）

#### anser服务目录结构 ####

```
  root |-- [service_name].proto
       |-- [service_name]_client.go
       |-- Makefile
       |-- service/
                    |-- main.go
                    |-- [service_name].go
                    |-- ...

```

将client放在服务顶层是为了方便其他模块使用，可以 

```

import (
    ".../service1"
    ".../service2"
)

```

而不是

```
import (
    service1 ".../service1/client"
    service2 ".../service2/client"
)

```

### anser组件 ###

* anser-config: anser配置，有哪些服务？服务有哪些接口？依赖关系？ 服务状态？
* http-gateway: 对外提供http访问接口
* grpc-gateway: 对外提供grpc访问接口
* anser-dashboard: 查询集群、服务和接口的状态