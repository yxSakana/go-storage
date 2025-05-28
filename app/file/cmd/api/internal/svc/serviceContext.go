package svc

import (
	"github.com/panjf2000/ants/v2"
	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/stores/redis"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"go-storage/app/file/cmd/api/internal/config"
	"go-storage/app/file/model"
)

const (
	_saveChunkPoolSize = 200
	_mergePoolSize     = 20
)

type ServiceContext struct {
	Config config.Config

	Redis *redis.Redis

	FileMetaModel    model.FileMetaModel
	UserFileRelModel model.UserFileRelModel
	UploadManager    UploadManager
	// kafka
	KPusherClient *kq.Pusher
	//KPusherClient *kafka.Writer
	SaveChunkPool *ants.Pool
	MergePool     *ants.Pool
}

func NewServiceContext(c config.Config) *ServiceContext {
	r := redis.MustNewRedis(c.RedisConf)
	p, err := ants.NewPool(_saveChunkPoolSize)
	if err != nil {
		panic(err)
	}
	p2, err := ants.NewPool(_mergePoolSize)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config: c,

		Redis: r,

		FileMetaModel:    model.NewFileMetaModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UserFileRelModel: model.NewUserFileRelModel(sqlx.NewMysql(c.DB.DataSource), c.Cache),
		UploadManager:    NewUploadManager(r),
		KPusherClient:    kq.NewPusher(c.KqPusherConf.Brokers, c.KqPusherConf.Topic),
		//KPusherClient: &kafka.Writer{
		//	Addr:        kafka.TCP(""),
		//	Topic:       "",
		//	Balancer:    &kafka.LeastBytes{},
		//	Compression: kafka.Snappy,
		//},
		SaveChunkPool: p,
		MergePool:     p2,
	}
}
