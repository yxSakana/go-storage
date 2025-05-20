package fileUpoad

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/file/cmd/api/internal/logic/fileUpoad"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"
	"go-storage/pkg/response"
)

// 文件元信息初始化
func InitFileMetaHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.InitFileMetaReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := fileUpoad.NewInitFileMetaLogic(r.Context(), svcCtx)
		resp, err := l.InitFileMeta(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
