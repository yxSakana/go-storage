package fileMeta

import (
	"context"

	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"

	"github.com/zeromicro/go-zero/core/logx"
)

type ConfigurateFileMetaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件元信息设置
func NewConfigurateFileMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ConfigurateFileMetaLogic {
	return &ConfigurateFileMetaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ConfigurateFileMetaLogic) ConfigurateFileMeta(req *types.ConfigurateFileMetaReq) (resp *types.ConfigurateFileMetaResp, err error) {
	// todo: add your logic here and delete this line

	return
}
