package user

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/app/user/cmd/api/internal/logic/user"
	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"
	"go-storage/pkg/response"
)

// register
func RegisterHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.RegisterReq
		if err := httpx.Parse(r, &req); err != nil {
			response.ParamError(r.Context(), w, err)
			return
		}

		l := user.NewRegisterLogic(r.Context(), svcCtx)
		resp, err := l.Register(&req)
		response.HttpResult(r.Context(), w, resp, err)
	}
}
