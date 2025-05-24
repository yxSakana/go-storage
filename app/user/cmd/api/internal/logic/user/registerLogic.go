package user

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"
	"go-storage/app/user/cmd/rpc/user"
)

type RegisterLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// register
func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *RegisterLogic) Register(req *types.RegisterReq) (resp *types.RegisterResp, err error) {
	_, err = l.svcCtx.UserRpc.Register(l.ctx, &user.RegisterReq{
		Email:    req.Email,
		Password: req.Password,
	})

	if err != nil {
		return nil, fmt.Errorf("register: %w", err)
	}

	return &types.RegisterResp{}, nil
}
