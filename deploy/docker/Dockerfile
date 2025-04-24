# 常用基础镜像:
# buildpack-deps maven gradle
# nginx httpd
# mysql redis postgres mongo
# alpine debian ubunta busybox
# golang python openjdk rust ruby ...

# or loads/alpine:3.8
FROM golang:latest

ENV GOPROXY=https://goproxy.cn,direct

ENV WORKDIR="/app"
WORKDIR ${WORKDIR}
ADD . ${WORKDIR}/

RUN go install github.com/cortesi/modd/cmd/modd@latest && go install github.com/air-verse/air@latest && go mod tidy

ENTRYPOINT ["modd"]
