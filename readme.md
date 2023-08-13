## 基于Go语言的仿云盘demo

### 项目结构说明

│  go.mod		: go module管理依赖文件

│  go.sum		: go module管理生成的版本管理文件

│  readme.md	: 本文件，项目说明

│  tree.md		: 记录文件目录结构树

├─cache		: 提供redis缓存支持，主要用于文件分块传输

│  └─redis

├─common		: 统一错误码(实际并没有怎么用)

├─config		: 统一项目配置(数据库、oss、mq、redis、监听地址等)

├─db			: 提供数据库(dao层)支持，包括创建连接池、提供crud接口

│  └─mysql

├─doc			: 项目相关，数据库表等

├─handler		: **原生Go语言模式** 下的handler方法

│  └─Gin-handler	: **Gin框架模式** 下的handler方法

├─meta			: 提供文件元信息结构和相关方法

├─mq			: 提供rabiitmq支持，添加、消费相关的接口

├─route			: **Gin框架模式** 下的路由-handler映射

├─service		: 项目启动入口

│  ├─Gin		: **Gin框架模式** 启动入口

│  ├─normal		: **原生Go语言模式** 的两个服务(云盘webApp、文件转存)启动入口

│  └─Microservice	: **微服务模式** 的主要实现(目前未完成 - 2019/7/28)

│      ├─account		: 用户相关微服务实现

│      └─apigw		: API网关

├─static		: 包含项目静态资源（页面、css、js等）

├─store		: 用于提供第三方文件云存储支持，目前只有阿里云oss

├─test		: 包含用于测试接口功能的简单脚本

└─util		: 包含用于hash加密、json转换的工具函数与结构

> - [ ] redis客户端切换成go-redis,项目用的那个太老了

<a href="https://www.freepik.com/free-vector/cute-astronaut-super-flying-cartoon-illustration_15644423.htm#query=astronaut&position=32&from_view=keyword&track=sph">Image by catalyststuff</a> on Freepik

### 项目启动

- rabiitmq转存任务启动

  > go run service/normal/transfer/main.go

- 原生Go语言模式云盘服务端启动

  > go run service/normal/upload/main.go

- Gin框架模式云盘服务端启动

  > go run service/Gin/main.go

- 微服务模式

  - 账户相关微服务启动

    > go run service/Microservice/account/main.go --registry=consul

  - API网关服务启动

    > go run service/Microservice/apigw/main.go --registry=consul

