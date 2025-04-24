package fileMeta

import (
	"context"
	"errors"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/app/file/model"
	"go-storage/app/token"
	"go-storage/pkg/gserr"
)

type QuickTransmissionLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件秒传通过hash值判断文件是否存在实现
func NewQuickTransmissionLogic(ctx context.Context, svcCtx *svc.ServiceContext) *QuickTransmissionLogic {
	return &QuickTransmissionLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *QuickTransmissionLogic) QuickTransmission(req *types.QuickTransmissionReq) (resp *types.QuickTransmissionResp, err error) {
	userId, err := token.GetUserId(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: %w", gserr.ErrServerCommon, err)
	}

	meta, err := l.svcCtx.FileMetaModel.FindOneByHash(l.ctx, req.Hash)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("%w: %w", gserr.ErrServerCommon, err)
	}

	// 文件记录不存在, 直接返回错误码，让前端上传去上传文件
	if meta == nil {
		return nil, gserr.ErrFileUpload
	}

	// 已经有这个文件记录, 创建文件与用户的映射关系
	_, err = l.svcCtx.UserFileRelModel.Insert(l.ctx, &model.UserFileRel{
		UserId:        userId,
		FileId:        meta.Id,
		FilenameAlias: req.FilenameAlias,
	})
	if err != nil {
		return nil, fmt.Errorf("%w: %w", gserr.ErrServerCommon, err)
	}

	return &types.QuickTransmissionResp{}, nil
}
