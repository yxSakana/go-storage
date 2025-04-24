package fileDownload

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/file/cmd/api/internal/logic/fileDownload"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/pkg/response"
)

func DownloadRangeHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.DownloadRangeReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := fileDownload.NewDownloadRangeLogic(r.Context(), svcCtx)
		_, err := l.DownloadRange(&req, w, r)
		if err != nil {
			response.HttpResult(r.Context(), w, nil, err)
		}
	}
}
