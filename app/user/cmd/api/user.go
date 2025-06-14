package main

import (
	"flag"
	"fmt"

	"go-storage/app/user/cmd/api/internal/config"
	"go-storage/app/user/cmd/api/internal/handler"
	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/pkg/response"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "etc/user.yaml", "the config file")

func main() {
	flag.Parse()

	var c config.Config
	conf.MustLoad(*configFile, &c)

	server := rest.MustNewServer(c.RestConf, rest.WithUnauthorizedCallback(response.UnauthorizedCallback))
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	handler.RegisterHandlers(server, ctx)

	fmt.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}
