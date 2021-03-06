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

### 工作进度纪录

#### 2019/7/10

Go云盘项目创建

#### 2019/7/11

后台文件上传下载接口编写与测试

#### 2019/7/13

后端用户模块api编写，实现登录与注册功能

#### 2019/7/14

用户文件表和api的编写，实现用户与文件的关联

#### 2019/7/15

秒传功能实现，通过hash值比较，判断文件是否已经在服务端存在

#### 2019/7/16

分块上传功能的部分实现，通过redis缓存记录各个分块上传的情况

#### 2019/7/18

阿里云oss云储存功能实现，在上传文件接口中，将文件上传到阿里云oss存储，下载时返回url，直接从oss下载

#### 2019/7/20

添加基于rabbitmq的异步文件转存功能，将程序拆分成两个模块 -- 文件上传 + 文件转存

#### 2019/7/23

将接口转换成Gin框架风格

#### 2019/7/24

将项目注册接口初步转为微服务架构

#### 2019/7/28

将项目转换为go module管理

实现项目的初步微服务化

添加了项目结构的说明

#### 2019/7/31

项目所有模块微服务化（未测试）

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

