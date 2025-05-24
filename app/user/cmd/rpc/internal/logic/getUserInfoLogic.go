package logic

import (
	"context"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/user/cmd/rpc/internal/svc"
	"go-storage/app/user/cmd/rpc/pb"
)

type GetUserInfoLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGetUserInfoLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GetUserInfoLogic {
	return &GetUserInfoLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GetUserInfoLogic) GetUserInfo(in *pb.GetUserInfoReq) (*pb.GetUserInfoResp, error) {
	// query by uid
	user, err := l.svcCtx.UserModel.FindOne(l.ctx, in.Id)
	if err != nil || user == nil {
		return nil, fmt.Errorf("get user(uid: %d) info: %w", in.Id, err)
	}

	return &pb.GetUserInfoResp{
		Userinfo: &pb.User{
			Id:       user.Id,
			Username: user.Username,
			Email:    user.Email,
			Avatar:   user.Avatar,
		},
	}, nil
}
