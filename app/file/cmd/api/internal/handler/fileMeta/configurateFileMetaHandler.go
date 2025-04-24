package fileMeta

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/file/cmd/api/internal/logic/fileMeta"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"

	"go-storage/pkg/response"
)

// 文件元信息设置
func ConfigurateFileMetaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ConfigurateFileMetaReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := fileMeta.NewConfigurateFileMetaLogic(r.Context(), svcCtx)
		resp, err := l.ConfigurateFileMeta(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
