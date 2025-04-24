package config

import (
	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/rest"
)

type Config struct {
	rest.RestConf

	JwtAuth struct {
		AccessSecret string
	}

	Redis redis.RedisConf

	DB struct {
		DataSource string
	}
	Cache cache.CacheConf
}
