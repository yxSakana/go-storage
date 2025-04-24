package middleware

import (
	"go-storage/pkg/gserr"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"

	"go-storage/pkg/response"
)

func RecoverMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				//logx.Error(fmt.Sprintf("%v\n%s", err, debug.Stack()))
				logx.Error(err)
				httpx.WriteJson(w, http.StatusInternalServerError, response.Response{
					Code: gserr.PanicError,
					Msg:  "Server Internal Error",
					Data: nil,
				})
			}
		}()

		next(w, r)
	}
}
