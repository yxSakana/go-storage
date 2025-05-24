package user

import (
	"context"
	"fmt"

	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"
	"go-storage/app/user/cmd/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type LoginLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// login
func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *LoginLogic) Login(req *types.LoginReq) (resp *types.LoginResp, err error) {
	loginResp, err := l.svcCtx.UserRpc.Login(l.ctx, &user.LoginReq{
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		return nil, fmt.Errorf("login error: %w", err)
	}

	return &types.LoginResp{
		UserId:       loginResp.UserId,
		Token:        loginResp.Token,
		ExpireAfter:  loginResp.ExpireAfter,
		RefreshAfter: loginResp.RefreshAfter,
	}, nil
}
