package fileUpoad

import (
	"go-storage/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-storage/app/file/cmd/api/internal/logic/fileUpoad"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
)

func CompleteUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.CompleteUploadReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := fileUpoad.NewCompleteUploadLogic(r.Context(), svcCtx)
		resp, err := l.CompleteUpload(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
