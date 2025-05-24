package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/user/cmd/api/internal/logic/user"
	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"

	"go-storage/pkg/response"
)

// activate
func ActivateHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.ActivateReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := user.NewActivateLogic(r.Context(), svcCtx)
		resp, err := l.Activate(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
