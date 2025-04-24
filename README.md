# Go-Storage

## Port

- user-api
  - 5001
  - 4001
- user-rpc
  - 5002
  - 4002
- file-api
  - 5003
  - 4003
- mysql
  - 3306
- redis
  - 6379
- gitlab
  - 8929  -- ui
  - 2222:22
  - 80
  - 22
- jenkins
  - 8080  -- ui
  - 50000
- elasticsearch
  - 9200
  - 9300
- kibana
  - 5601  -- ui
- [jaeger](https://www.jaegertracing.io/docs/2.5/apis/)
  - 16686  -- ui
  - 4317
  - 14268
  - 4318
  - 5778
  - 9411
- prometheus
  - 9090 -- ui
- grafana
  - 3000 -- ui

## Feature

- [x] 账号管理
- [x] 文件秒传
- [x] 断点续传
- [x] 分块上传
- [x] 文件下载
- [x] 微服务架构
- [x] 支持大规模扩展(容器化、服务自治、监控)

## Run

```shell
docker compose -f ./deploy/docker/compose.yaml up -d
```

## 技术栈

- gateway
  - nginx
- [x] frame
  - go-zero
- [x] call
  - 外部: RESTful  API
  - 内部: gRPC
- 消息队列
  - kafka
- [x] DB
  - MySQL
  - Redis
- FileStorage
  - MinIO
- 日志系统
  - [x] Filebeat
  - Logstash
  - kafka
  - [x] elasticsearch
  - [x] kibana
- 系统监控
  - [x] Prometheus
  - Grafana
- 链路追踪
  - Jaeger
- CI
  - gitlab
  - Jenkins
- CD
  - [x] docker
  - harbor
  - wayne
  - k8s
- 自动监听文件改动并编译重启
  - [x] air
  - [x] modd
