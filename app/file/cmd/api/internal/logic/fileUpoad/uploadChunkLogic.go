package fileUpoad

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/pkg/file"
	"go-storage/pkg/gserr"
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
	if !l.svcCtx.UploadManager.Exists(req.FileHash) {
		return nil, fmt.Errorf("file %s not exists", req.FileHash)
	}
	if l.svcCtx.UploadManager.IsChunkCompleted(l.ctx, req.FileHash, req.ChunkIndex) {
		return &types.UploadChunkResp{}, nil
	}
	// Save chunk file
	chunkDir := fmt.Sprintf(types.ChunkDirf, req.FileHash)
	chunkFilename := fmt.Sprintf("%s/%d.chunk", chunkDir, req.ChunkIndex)
	err = file.SaveFileHeader(req.ChunkFileHeader, chunkFilename)
	if err != nil {
		return nil, fmt.Errorf("%w: write file %s: %v", gserr.ErrServerCommon, chunkFilename, err)
	}
	// verify file hash
	chunkHash, err := file.CalculateHash(chunkFilename, "md5")
	if err != nil {
		return nil, fmt.Errorf("%w: calculate file hash: %v", gserr.ErrServerCommon, err)
	}
	if chunkHash != req.ChunkHash {
		return nil, fmt.Errorf("%w: file hash %s is not equal to expected hash %s",
			gserr.ErrFileHashIncomplete, chunkHash, req.ChunkHash)
	}

	// Set flag in Redis
	if err := l.svcCtx.UploadManager.CompletedChunk(l.ctx, req.FileHash, req.ChunkIndex); err != nil {
		return nil, err
	}

	return &types.UploadChunkResp{}, nil
}
