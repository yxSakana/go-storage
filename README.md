# Go-Storage

## Port

- user-api: 5001 4001
- user-rpc: 5002 4002
- file-api: 5003 4003
- mysql: 3306
- redis: 6379
- elasticsearch: 9200 9300
- kibana: 5601   -- ui
- jaeger: 16686 4317 14268 4318 5778 9411 -- ui
- prometheus: 9090 -- ui
- grafana: 3000 -- ui

## Feature

- [x] 账号管理
- [x] 文件秒传
- [x] 断点续传
- [x] 分块上传
- [x] 文件下载
- [x] 微服务架构
- [x] 支持大规模扩展(容器化、服务自治、监控)

## 技术栈

- gateway
  - nginx
- frame
  - go-zero
- call
  - 外部: RESTful  API
  - 内部: gRPC
- 消息队列
  - kafka
- DB
  - MySQL
  - Redis
- FileStorage
  - MinIO
- 日志系统
  - Filebeat
  - Logstash
  - kafka
  - elasticsearch
  - kibana
- 系统监控
  - Prometheus
  - Grafana
- 链路追踪
  - Jaeger
- CI
  - gitlab
  - Jenkins
- CD
  - docker
  - harbor
  - wayne
  - k8s
- 自动监听文件改动并编译重启
  - air
  - modd
