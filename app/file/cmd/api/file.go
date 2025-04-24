package main

import (
	"flag"
	"fmt"

	"go-storage/app/file/cmd/api/internal/config"
	"go-storage/app/file/cmd/api/internal/handler"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/middleware"
	"go-storage/pkg/response"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/file.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(response.UnauthorizedCallback))
	defer server.Stop()
	server.Use(middleware.RecoverMiddleware)

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
