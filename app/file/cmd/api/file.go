package main

import (
	"context"
	"flag"
	"fmt"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/core/service"
	"github.com/zeromicro/go-zero/rest"

	"go-storage/app/file/cmd/api/internal/config"
	"go-storage/app/file/cmd/api/internal/handler"
	"go-storage/app/file/cmd/api/internal/mqs"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/middleware"
	"go-storage/pkg/response"
)

var configFile = flag.String("f", "etc/file.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c, conf.UseEnv())

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(response.UnauthorizedCallback))
	defer server.Stop()
	server.Use(middleware.RecoverMiddleware)

	svcCtx := svc.NewServiceContext(c)
	ctx := context.Background()
	handler.RegisterHandlers(server, svcCtx)
	serviceGroup := service.NewServiceGroup()
	defer serviceGroup.Stop()
	serviceGroup.Add(server)

	for _, mq := range mqs.Consumers(c, ctx, svcCtx) {
		serviceGroup.Add(mq)
	}

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	serviceGroup.Start()
	//server.Start()
}
