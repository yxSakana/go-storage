package fileUpoad

import (
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/file/cmd/api/internal/logic/fileUpoad"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/pkg/response"
)

func UploadChunkHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var req types.UploadChunkReq
		var logicInput types.UploadChunkInput
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}
		files := r.MultipartForm.File

		if chunkFiles, ok := files["chunk_data"]; ok {
			if len(chunkFiles) != 1 {
				response.ParamError(r.Context(), w, fmt.Errorf("number of parameter 'chunk_data' is incorrect"))
				return
			}
			logicInput = types.UploadChunkInput{
				UploadChunkReq:  req,
				ChunkFileHeader: chunkFiles[0],
			}
		} else {
			response.ParamError(r.Context(), w, fmt.Errorf("not has param: chunk_data"))
			return
		}

		l := fileUpoad.NewUploadChunkLogic(r.Context(), svcCtx)
		resp, err := l.UploadChunk(&logicInput)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
