package types

import (
	"mime/multipart"
)

const (
	CacheMetaKeyf  string = "gsCache:upload:%s:meta"
	CacheChunkKeyf string = "gsCache:upload:%s:chunk"
	ChunkDirf      string = "/tmp/uploads/chunk/%s"
)

type UploadChunkInput struct {
	UploadChunkReq
	ChunkFileHeader *multipart.FileHeader
}
