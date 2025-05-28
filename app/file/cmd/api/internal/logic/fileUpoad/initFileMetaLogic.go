package fileUpoad

import (
	"context"
	"errors"
	"fmt"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/app/file/model"
	"go-storage/pkg/gserr"
)

type InitFileMetaLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// 文件元信息初始化
func NewInitFileMetaLogic(ctx context.Context, svcCtx *svc.ServiceContext) *InitFileMetaLogic {
	return &InitFileMetaLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *InitFileMetaLogic) InitFileMeta(req *types.InitFileMetaReq) (resp *types.InitFileMetaResp, err error) {
	meta, err := l.svcCtx.FileMetaModel.FindOneByHash(l.ctx, req.Hash)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("%w: find by hash: %w", gserr.ErrServerCommon, err)
	}
	if meta != nil {
		return nil, fmt.Errorf("%w: init file meate: hash(%s) exists", gserr.ErrServerCommon, req.Hash)
	}

	if l.svcCtx.UploadManager.Exists(req.Hash) {
		meta, err := l.svcCtx.UploadManager.GetMeta(l.ctx, req.Hash)
		if err == nil {
			return &types.InitFileMetaResp{
				Hash:       meta.FileHash,
				ChunkSize:  meta.ChunkSize,
				ChunkCount: meta.ChunkCount,
			}, nil
		}
	}

	us := &svc.UploadInfo{
		FileHash:   req.Hash,
		FileSize:   req.FileSize,
		ChunkSize:  req.ChunkSize,
		ChunkCount: int((req.FileSize + req.ChunkSize - 1) / req.ChunkSize),
		FilePath:   fmt.Sprintf("data/uploads/%s%s", req.Hash, filepath.Ext(req.FilenameAlias)),
	}
	resp = &types.InitFileMetaResp{
		Hash:       req.Hash,
		ChunkSize:  req.ChunkSize,
		ChunkCount: us.ChunkCount,
	}
	// file meta to redis
	if err := l.svcCtx.UploadManager.InitMeta(l.ctx, *us); err != nil {
		return nil, err
	}

	return resp, nil
}
