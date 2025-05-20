package types

import (
	"mime/multipart"
)

const (
	ChunkDirf string = "/tmp/uploads/chunk/%s"
)

type UploadChunkInput struct {
	UploadChunkReq
	ChunkFileHeader *multipart.FileHeader
}
