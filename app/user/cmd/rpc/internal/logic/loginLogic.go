package logic

import (
	"context"
	"fmt"
	"go-storage/pkg/crypto"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/user/cmd/rpc/internal/svc"
	"go-storage/app/user/cmd/rpc/pb"
)

type LoginLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewLoginLogic(ctx context.Context, svcCtx *svc.ServiceContext) *LoginLogic {
	return &LoginLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *LoginLogic) Login(in *pb.LoginReq) (*pb.LoginResp, error) {
	// query Mobile
	user, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil || user == nil {
		return nil, fmt.Errorf("login: %w", err)
	}
	// confirm password
	if ok := crypto.CheckPassword(in.Password, user.Password); !ok {
		return nil, fmt.Errorf("login: invalid password")
	}
	// generate token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{UserId: user.Id})
	if err != nil {
		return nil, fmt.Errorf("login failed: generate token failed: %w", err)
	}

	return &pb.LoginResp{
		UserId:       user.Id,
		Token:        tokenResp.Token,
		ExpireAfter:  tokenResp.ExpireAfter,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
