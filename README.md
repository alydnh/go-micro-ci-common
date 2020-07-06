# go-micro-ci-common
基于Golang开发的微CI公共库

## 2020-07-03 日志模块新增 v0.0.2-logs
1. call 与 then 支持自定义函数调用
2. 原 call 与 then 重构为 handle 与 thenHandle
3. 新增 LogrusScopeWriter 与 io.Writer的适合器

## 2020-07-02 新增日志模块 v0.0.1-logs
1. micro_logrus.go go-micro日志适配器
2. logrus_scope.go 基于logrus entry的封装，支持then和onerror调用，

## 2020-06-29 新增registry定义
```yaml
registry:
  type: consul
  address: "consul"
  port: 8500
  useSSL: false
```
daemon程序应该会根据registry定义，连接服务注册表，如未能连接，则偿式从yaml找寻相应的服务定义，启动此服务。
## 2020-06-12 初始提交
1. utils目录 常用公共函数
2. yaml 基于Yaml的ci基本结构体定义
