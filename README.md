# go-micro-ci-common
基于Golang开发的微CI公共库

## 2020-07-10 logrusScope 新增 catch 方法
1. 可以在call与then这后调用catch
2. 返回错误则会继续调用onError, 如返回空，继续then

## 2020-07-10 service 新增 disableCommonEnv
1. 部署时，如disableCommonEnv:true时，容器启动不会使用commonEnv的定义
2. 容器环境变量覆盖优先级 commonEnv(if not 如disableCommonEnv) < registryEnv < serviceEnv

## 2020-07-06 logrusScope.internal call panic策略
1. 在call之前进行defer recover
2. 上报panic的stack
3. 将panic转换成error

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
