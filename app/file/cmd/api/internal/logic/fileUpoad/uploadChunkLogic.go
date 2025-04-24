package fileUpoad

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/pkg/file"
	"go-storage/pkg/gserr"
	"strconv"
)

type UploadChunkLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewUploadChunkLogic(ctx context.Context, svcCtx *svc.ServiceContext) *UploadChunkLogic {
	return &UploadChunkLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *UploadChunkLogic) UploadChunk(req *types.UploadChunkInput) (resp *types.UploadChunkResp, err error) {
	key := fmt.Sprintf(types.CacheChunkKeyf, req.FileId)
	if v, err := l.svcCtx.Redis.HgetCtx(l.ctx, key, strconv.Itoa(req.ChunkIndex)); err != nil {
		return nil, fmt.Errorf("%w: %w", gserr.ErrServerCommon, err)
	} else {
		if v == "1" {
			return &types.UploadChunkResp{}, nil
		}
	}
	// Save chunk file
	chunkDir := fmt.Sprintf(types.ChunkDirf, req.FileId)
	chunkFilename := fmt.Sprintf("%s/%d.chunk", chunkDir, req.ChunkIndex)
	err = file.SaveFileHeader(req.ChunkFileHeader, chunkFilename)
	if err != nil {
		return nil, fmt.Errorf("%w: write file %s: %v", gserr.ErrServerCommon, chunkFilename, err)
	}

	// Set flag in Redis
	err = l.svcCtx.Redis.HsetCtx(l.ctx, key, strconv.Itoa(req.ChunkIndex), "1")
	if err != nil {
		return nil, fmt.Errorf("%w: %w", gserr.ErrServerCommon, err)
	}

	return &types.UploadChunkResp{}, nil
}
