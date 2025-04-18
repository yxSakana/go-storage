package user

import (
	"go-storage/pkg/response"
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
	"go-storage/app/user/cmd/api/internal/logic/user"
	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"
)

// get user info
func DetailHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.GetUserInfoReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := user.NewDetailLogic(r.Context(), svcCtx)
		resp, err := l.Detail(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
