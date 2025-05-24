package user

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/user/cmd/api/internal/svc"
	"go-storage/app/user/cmd/api/internal/types"
	"go-storage/app/user/cmd/rpc/user"
)

type DetailLogic struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

// get user info
func NewDetailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *DetailLogic {
	return &DetailLogic{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *DetailLogic) Detail(req *types.GetUserInfoReq) (resp *types.GetUserInfoResp, err error) {
	userInfoResp, err := l.svcCtx.UserRpc.GetUserInfo(l.ctx, &user.GetUserInfoReq{
		Id: req.UserId,
	})
	if err != nil {
		return nil, fmt.Errorf("get user info(%d): %w", req.UserId, err)
	}

	return &types.GetUserInfoResp{
		UserInfo: types.User{
			Id:       userInfoResp.Userinfo.Id,
			Username: userInfoResp.Userinfo.Username,
			Email:    userInfoResp.Userinfo.Email,
			Avatar:   userInfoResp.Userinfo.Avatar,
		},
	}, nil
}
