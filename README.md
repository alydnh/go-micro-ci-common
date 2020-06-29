# go-micro-ci-common
基于Golang开发的微CI公共库

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
