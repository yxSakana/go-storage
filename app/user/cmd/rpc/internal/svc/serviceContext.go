package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"github.com/zeromicro/go-zero/zrpc"

	"go-storage/app/email/cmd/rpc/email"
	"go-storage/app/user/cmd/rpc/internal/config"
	"go-storage/app/user/model"
)

type ServiceContext struct {
	Config config.Config

	Redis *redis.Redis

	UserModel model.UserModel

	EmailRpc      email.Email
	RegisterCache *RegisterCache
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := redis.MustNewRedis(c.RedisConf)
	return &ServiceContext{
		Config: c,

		Redis: r,

		UserModel:     model.NewUserModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		EmailRpc:      email.NewEmail(zrpc.MustNewClient(c.EmailRpcConf)),
		RegisterCache: NewRegisterCache(r),
	}
}
