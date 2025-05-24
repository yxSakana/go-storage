package user

import (
	"context"

	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"
	"go-storage/app/user/cmd/rpc/user"

	"github.com/zeromicro/go-zero/core/logx"
)

type ActivateLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// activate
func NewActivateLogic(ctx context.Context, svcCtx *svc.ServiceContext) *ActivateLogic {
	return &ActivateLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *ActivateLogic) Activate(req *types.ActivateReq) (resp *types.ActivateResp, err error) {
	_, err = l.svcCtx.UserRpc.ActivateAccount(l.ctx, &user.ActivateAccountReq{VerifyToken: req.VerifyToken})

	return &types.ActivateResp{}, err
}
