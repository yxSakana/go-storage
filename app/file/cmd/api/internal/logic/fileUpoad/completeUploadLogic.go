package fileUpoad

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/core/stores/sqlx"

	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/app/file/model"
	"go-storage/app/token"
	"go-storage/pkg/file"
	"go-storage/pkg/gserr"
)

type MergeMessage struct {
	FileHash string
	UserId   uint64
}

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

	fm, err := l.svcCtx.FileMetaModel.FindOneByHash(l.ctx, req.FileHash)
	if err == nil {
		ret, err := l.svcCtx.UserFileRelModel.FindOneByUserIdFileId(l.ctx, userId, fm.Id)
		if err == nil && ret != nil {
			return &types.CompleteUploadResp{}, nil
		}
	}

	if !l.svcCtx.UploadManager.Exists(req.FileHash) {
		return nil, fmt.Errorf("%w: %s", gserr.ErrFileMetaUninitialized, req.FileHash)
	}

	// check the chunk flag saved in Redis
	if !l.svcCtx.UploadManager.CanMerge(l.ctx, req.FileHash) {
		return nil, fmt.Errorf("%w: file_hash: %s",
			gserr.ErrFileIncompleteChunk, req.FileHash)
	}

	// kafka write
	mess, err := json.Marshal(MergeMessage{
		FileHash: req.FileHash,
		UserId:   userId,
	})
	if err != nil {
		return nil, err
	}
	m := string(mess)
	if err := l.svcCtx.KPusherClient.Push(l.ctx, m); err != nil {
		return nil, err
	}
	return &types.CompleteUploadResp{}, nil
}

func (l *CompleteUploadLogic) MergeChunks(meta *svc.UploadInfo, userId uint64) error {
	// 合并文件 && 验证hash
	chunkDir := fmt.Sprintf(types.ChunkDirf, meta.FileHash)
	filePath := meta.FilePath
	if err := mergeChunks(meta.ChunkCount, chunkDir, filePath); err != nil {
		return err
	}
	if err := file.VerifyFileHash(filePath, meta.FileHash); err != nil {
		return err
	}
	// Save to MySQL
	if err := l.svcCtx.FileMetaModel.Trans(l.ctx, func(context context.Context, session sqlx.Session) error {

		fileMetaInsertRet, err := l.svcCtx.FileMetaModel.InsertWithSession(l.ctx, session, &model.FileMeta{
			Hash:   meta.FileHash,
			Size:   int64(meta.FileSize),
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
			FilenameAlias: meta.Filename,
		}); err != nil {
			return fmt.Errorf("%w: insert file rel model: %w", gserr.ErrServerCommon, err)
		}

		return nil
	}); err != nil {
		return err
	}

	_ = l.svcCtx.UploadManager.CompletedMerge(l.ctx, meta.FileHash)
	return nil
}

func mergeChunks(chunkCount int, chunkDir string, outputPath string) error {
	saveDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(saveDir, os.ModePerm); err != nil {
		return err
	}
	f, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer f.Close()

	for i := 0; i < chunkCount; i++ {
		chunkPath := fmt.Sprintf("%s/%d.chunk", chunkDir, i)
		chunkData, err := os.ReadFile(chunkPath)
		if err != nil {
			return err
		}
		if _, err := f.Write(chunkData); err != nil {
			return err
		}
	}
	return nil
}
