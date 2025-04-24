package fileMeta

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/file/cmd/api/internal/logic/fileMeta"
	"go-storage/app/file/cmd/api/internal/svc"
	"go-storage/app/file/cmd/api/internal/types"

	"go-storage/pkg/response"
)

// 文件秒传通过hash值判断文件是否存在实现
func QuickTransmissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.QuickTransmissionReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := fileMeta.NewQuickTransmissionLogic(r.Context(), svcCtx)
		resp, err := l.QuickTransmission(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
