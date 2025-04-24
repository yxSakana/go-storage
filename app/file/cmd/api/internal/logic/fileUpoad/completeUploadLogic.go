package fileUpoad

import (
	"context"
	"fmt"
	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/app/file/model"
	"go-storage/app/token"
	"go-storage/pkg/gserr"
	"os"
	"path/filepath"
	"strconv"
)

type CompleteUploadLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewCompleteUploadLogic(ctx context.Context, svcCtx *svc.ServiceContext) *CompleteUploadLogic {
	return &CompleteUploadLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *CompleteUploadLogic) CompleteUpload(req *types.CompleteUploadReq) (resp *types.CompleteUploadResp, err error) {
	userId, err := token.GetUserId(l.ctx)
	if err != nil {
		return nil, fmt.Errorf("%w: get user id: %w", gserr.ErrServerCommon, err)
	}
	fileIdUint64, err := strconv.ParseUint(req.FileId, 10, 64)
	// check DB 如果存在直接返回
	if err == nil {
		ret, err := l.svcCtx.UserFileRelModel.FindOneByUserIdFileId(l.ctx, userId, fileIdUint64)
		if err == nil && ret != nil {
			return &types.CompleteUploadResp{}, nil
		}
	}
	// 检查key是否存在
	metaKey := fmt.Sprintf(types.CacheMetaKeyf, req.FileId)
	chunkKey := fmt.Sprintf(types.CacheChunkKeyf, req.FileId)

	if exist, err := l.svcCtx.Redis.Exists(metaKey); err != nil {
		return nil, fmt.Errorf("%w: get redis meta key: %w", gserr.ErrServerCommon, err)
	} else if !exist {
		return nil, fmt.Errorf("%w: %s", gserr.ErrFileMetaUninitialized, req.FileId)
	}
	// get meta from redis
	meta, err := l.svcCtx.Redis.HgetallCtx(l.ctx, metaKey)
	if err != nil {
		return nil, fmt.Errorf("%w: get meta: %w", gserr.ErrAttachedMsgError, err)
	}
	// check the chunk flag saved in Redis
	chunkFlags, err := l.svcCtx.Redis.HgetallCtx(l.ctx, chunkKey)
	if err != nil {
		return nil, fmt.Errorf("%w: get chunk flags: %w", gserr.ErrServerCommon, err)
	}
	chunkCount := len(chunkFlags)
	for i := 0; i < chunkCount; i++ {
		if chunkFlags[strconv.Itoa(i)] != "1" {
			return nil, fmt.Errorf("%w: file_id: %s, chunk_index: %d",
				gserr.ErrFileIncompleteChunk, req.FileId, i)
		}
	}

	// 合并文件
	// TODO: 验证hash
	chunkDir := fmt.Sprintf(types.ChunkDirf, req.FileId)
	filePath := meta["filePath"]
	saveDir := filepath.Dir(filePath)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return nil, fmt.Errorf("%w: create save dir: %w", gserr.ErrServerCommon, err)
	}
	f, err := os.Create(filePath)
	if err != nil {
		return nil, fmt.Errorf("%w: create file %s: %w", gserr.ErrServerCommon, req.FileId, err)
	}
	for i := 0; i < chunkCount; i++ {
		chunkPath := fmt.Sprintf("%s/%d.chunk", chunkDir, i)
		chunkDta, err := os.ReadFile(chunkPath)
		if err != nil {
			return nil, fmt.Errorf("%w: read file %s: %w", gserr.ErrServerCommon, req.FileId, err)
		}
		if _, err := f.Write(chunkDta); err != nil {
			return nil, fmt.Errorf("%w: write file %s: %w", gserr.ErrServerCommon, req.FileId, err)
		}
	}

	// Save to MySQL
	fileSize, err := strconv.ParseInt(meta["fileSize"], 10, 64)
	if err != nil {
		return nil, fmt.Errorf("%w: file size(%s): %w", gserr.ErrServerCommon, req.FileId, err)
	}
	if err := l.svcCtx.FileMetaModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {

		fileMetaInsertRet, err := l.svcCtx.FileMetaModel.InsertWithSession(l.ctx, session, &model.FileMeta{
			Hash:   meta["fileHash"],
			Size:   fileSize,
			Path:   filePath,
			Status: 1,
		})
		if err != nil {
			return fmt.Errorf("insert file meta: %w", err)
		}
		fileMetaId, err := fileMetaInsertRet.LastInsertId()
		if err != nil {
			return fmt.Errorf("%w: get file meta id: %w", gserr.ErrServerCommon, err)
		}

		// 文件 <-> 用户 关系映射
		if _, err := l.svcCtx.UserFileRelModel.InsertWithSession(l.ctx, session, &model.UserFileRel{
			UserId:        userId,
			FileId:        uint64(fileMetaId),
			FilenameAlias: meta["filename"],
		}); err != nil {
			return fmt.Errorf("%w: insert file rel model: %w", gserr.ErrServerCommon, err)
		}

		return nil
	}); err != nil {
		return nil, fmt.Errorf("%w: Database: %w", gserr.ErrServerCommon, err)
	}

	return &types.CompleteUploadResp{}, nil
}
