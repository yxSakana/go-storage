package svc

import (
	"github.com/zeromicro/go-zero/zrpc"
	"go-storage/app/user/cmd/api/internal/config"
	"go-storage/app/user/cmd/rpc/user"
)

type ServiceContext struct {
	Config  config.Config
	UserRpc user.UserZrpcClient
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:  c,
		UserRpc: user.NewUserZrpcClient(zrpc.MustNewClient(c.UserRpcConf)),
	}
}
