package mqs

import (
	"context"
	"go-storage/app/file/cmd/api/internal/logic/fileUpoad"

	"github.com/zeromicro/go-queue/kq"
	"github.com/zeromicro/go-zero/core/service"

	"go-storage/app/file/cmd/api/internal/config"
	"go-storage/app/file/cmd/api/internal/svc"
)

func Consumers(c config.Config, ctx context.Context, svcContext *svc.ServiceContext) []service.Service {

	return []service.Service{
		kq.MustNewQueue(c.KqConsumerConf, fileUpoad.NewMergeConsumer(ctx, svcContext)),
	}

}
