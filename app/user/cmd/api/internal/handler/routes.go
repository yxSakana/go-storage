// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.2

package handler

import (
	"net/http"

	user "go-storage/app/user/cmd/api/internal/handler/user"
	"go-storage/app/user/cmd/api/internal/svc"

	"github.com/zeromicro/go-zero/rest"
)

func RegisterHandlers(server *rest.Server, serverCtx *svc.ServiceContext) {
	server.AddRoutes(
		[]rest.Route{
			{
				// activate
				Method:  http.MethodGet,
				Path:    "/activate",
				Handler: user.ActivateHandler(serverCtx),
			},
			{
				// login
				Method:  http.MethodPost,
				Path:    "/login",
				Handler: user.LoginHandler(serverCtx),
			},
			{
				// register
				Method:  http.MethodPost,
				Path:    "/register",
				Handler: user.RegisterHandler(serverCtx),
			},
			{
				Method:  http.MethodGet,
				Path:    "/register/resend-email",
				Handler: user.ResendActivateEmailHandler(serverCtx),
			},
		},
		rest.WithPrefix("/api/v1/user"),
	)

	server.AddRoutes(
		[]rest.Route{
			{
				// get user info
				Method:  http.MethodGet,
				Path:    "/detail",
				Handler: user.DetailHandler(serverCtx),
			},
		},
		rest.WithJwt(serverCtx.Config.JwtAuth.AccessSecret),
		rest.WithPrefix("/api/v1/user"),
	)
}
