package svc

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/zeromicro/go-zero/core/stores/redis"
)

const (
	_cacheMetaKeyf string = "gsCache:upload:%s:meta"

	_expireSec int = 12 * 60 * 60
)

type UploadInfo struct {
	FileHash   string
	Filename   string
	FileSize   uint64
	ChunkSize  uint64
	ChunkCount int
	FilePath   string
	ChunkState []bool
}

func (us *UploadInfo) ToRedis() map[string]string {
	s := fmt.Sprintf("%t", us.ChunkState)
	return map[string]string{
		"fileHash":   us.FileHash,
		"filename":   us.Filename,
		"fileSize":   fmt.Sprintf("%d", us.FileSize),
		"chunkSize":  fmt.Sprintf("%d", us.ChunkSize),
		"chunkCount": fmt.Sprintf("%d", us.ChunkCount),
		"filePath":   us.FilePath,
		"chunkState": s[1 : len(s)-1],
	}
}

func FromRedis(r map[string]string) *UploadInfo {
	if r == nil {
		return nil
	}
	fileSize, _ := strconv.ParseUint(r["chunkSize"], 10, 64)
	chunkSize, _ := strconv.ParseUint(r["chunkSize"], 10, 64)
	chunkCount, _ := strconv.Atoi(r["chunkCount"])
	chunkState := make([]bool, chunkCount)
	for i := 0; i < chunkCount; i++ {
		chunkState[i] = r["chunkState"] == "true"
	}
	return &UploadInfo{
		FileHash:   r["fileHash"],
		Filename:   r["filename"],
		FileSize:   fileSize,
		ChunkSize:  chunkSize,
		ChunkCount: chunkCount,
		FilePath:   r["filePath"],
		ChunkState: chunkState,
	}
}

type UploadManager interface {
	InitMeta(context.Context, UploadInfo) error
	SetMeta(context.Context, UploadInfo) error
	GetMeta(ctx context.Context, hash string) (*UploadInfo, error)
	Exists(hash string) bool
	CompletedChunk(ctx context.Context, hash string, chunkIndex int) error
	CompletedMerge(ctx context.Context, hash string) error
	IsChunkCompleted(ctx context.Context, hash string, chunkIndex int) bool
	CanMerge(ctx context.Context, hash string) bool
}

type uploadManager struct {
	rdb *redis.Redis
}

func NewUploadManager(c *redis.Redis) UploadManager {
	return &uploadManager{rdb: c}
}

var _ UploadManager = (*uploadManager)(nil)

func (m *uploadManager) InitMeta(ctx context.Context, us UploadInfo) error {
	us.ChunkState = make([]bool, us.ChunkCount)
	return m.SetMeta(ctx, us) // todo:
}

func (m *uploadManager) SetMeta(ctx context.Context, us UploadInfo) (err error) {
	err = m.rdb.HmsetCtx(ctx, getMetaKey(us.FileHash), us.ToRedis())
	_ = m.rdb.Expire(getMetaKey(us.FileHash), _expireSec)
	return
}

func (m *uploadManager) GetMeta(ctx context.Context, hash string) (us *UploadInfo, err error) {
	meta, err := m.rdb.HgetallCtx(ctx, getMetaKey(hash))
	us = FromRedis(meta)
	return
}

func (m *uploadManager) Exists(hash string) bool {
	exist, err := m.rdb.Exists(getMetaKey(hash))
	return exist && err == nil
}

func (m *uploadManager) ExistsByHash(ctx context.Context, hash string) bool {
	exist, err := m.rdb.HmgetCtx(ctx, getMetaKey(hash), "fileHash")
	return exist != nil && err == nil
}

func (m *uploadManager) CompletedChunk(ctx context.Context, hash string, chunkIndex int) (err error) {
	results, err := m.rdb.HmgetCtx(ctx, getMetaKey(hash), "chunkState")
	if err != nil {
		return err
	}
	if len(results) != 1 {
		return fmt.Errorf("not found chunk state")
	}
	chunks := strings.Split(results[0], " ")
	chunks[chunkIndex] = "true"
	return m.rdb.HmsetCtx(ctx, getMetaKey(hash), map[string]string{
		"chunkState": strings.Join(chunks, " "),
	})
}

func (m *uploadManager) CompletedMerge(ctx context.Context, hash string) error {
	_, err := m.rdb.DelCtx(ctx, getMetaKey(hash))
	return err
}

func (m *uploadManager) IsChunkCompleted(ctx context.Context, hash string, chunkIndex int) bool {
	if v, err := m.rdb.HgetCtx(ctx, getMetaKey(hash), strconv.Itoa(chunkIndex)); err != nil {
		return false
	} else {
		return v == "true"
	}
}

func (m *uploadManager) CanMerge(ctx context.Context, hash string) bool {
	results, err := m.rdb.HmgetCtx(ctx, getMetaKey(hash), "chunkState")
	if err != nil {
		return false
	}

	chunks := strings.Split(results[0], " ")
	for _, chunk := range chunks {
		if chunk != "true" {
			return false
		}
	}
	return true
}

func getMetaKey(fid string) string {
	return fmt.Sprintf(_cacheMetaKeyf, fid)
}

//func getChunkKey(fid string) string {
//	return fmt.Sprintf(_cacheChunkKeyf, fid)
//}
