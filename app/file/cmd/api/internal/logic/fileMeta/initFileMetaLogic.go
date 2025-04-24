package fileMeta

import (
	"context"
	"path/filepath"

	//"errors"
	"fmt"
	"strconv"

	"github.com/cockroachdb/errors"
	"github.com/google/uuid"
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

	fileId := uuid.NewString()
	metaKey := fmt.Sprintf(types.CacheMetaKeyf, fileId)
	chunkKey := fmt.Sprintf(types.CacheChunkKeyf, fileId)
	chunkCount := (req.FileSize + req.ChunkSize - 1) / req.ChunkSize
	resp = &types.InitFileMetaResp{
		Id:         fileId,
		ChunkSize:  req.ChunkSize,
		ChunkCount: int(chunkCount),
	}
	// 检查Redis中是否已经存在
	if exist, err := l.svcCtx.Redis.Exists(metaKey); err != nil {
		return nil, err
	} else if exist {
		return resp, nil
	}
	// file meta to redis
	if err := l.svcCtx.Redis.HmsetCtx(l.ctx, metaKey, map[string]string{
		"fileHash":  req.Hash,
		"fileSize":  strconv.FormatUint(req.FileSize, 10),
		"chunkSize": strconv.FormatUint(req.ChunkSize, 10),
		"filePath":  fmt.Sprintf("data/uploads/%s%s", req.Hash, filepath.Ext(req.FilenameAlias)), // TODO: 文件后缀/类型, 使用mine映射
	}); err != nil {
		return nil, fmt.Errorf("%w: init file meate: save to cache: %w", gserr.ErrServerCommon, err)
	}
	_ = l.svcCtx.Redis.ExpireCtx(l.ctx, metaKey, 24*60*60)
	// file chunks to redis
	chunkMap := make(map[string]string, chunkCount)
	for i := range chunkCount {
		chunkMap[strconv.FormatUint(i, 10)] = "0"
	}
	if err := l.svcCtx.Redis.HmsetCtx(l.ctx, chunkKey, chunkMap); err != nil {
		return nil, fmt.Errorf("%w: init file meate: save to cache: %w", gserr.ErrServerCommon, err)
	}
	_ = l.svcCtx.Redis.ExpireCtx(l.ctx, chunkKey, 24*60*60)

	return resp, nil
}
