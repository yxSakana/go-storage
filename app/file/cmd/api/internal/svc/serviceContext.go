package svc

import (
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"go-storage/app/file/cmd/api/internal/config"
	"go-storage/app/file/model"
)

type ServiceContext struct {
	Config config.Config

	Redis *redis.Redis

	FileMetaModel    model.FileMetaModel
	UserFileRelModel model.UserFileRelModel
	UploadManager    UploadManager
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := redis.MustNewRedis(c.Redis)
	return &ServiceContext{
		Config: c,

		Redis: r,

		FileMetaModel:    model.NewFileMetaModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserFileRelModel: model.NewUserFileRelModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UploadManager:    NewUploadManager(r),
	}
}
