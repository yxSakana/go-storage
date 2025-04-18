package logic

import (
	"context"
	"errors"
	"fmt"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/user/cmd/rpc/internal/svc"
	"go-storage/app/user/cmd/rpc/pb"
	"go-storage/app/user/model"
	"go-storage/pkg/crypto"
	"go-storage/pkg/gserr"
)

type RegisterLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewRegisterLogic(ctx context.Context, svcCtx *svc.ServiceContext) *RegisterLogic {
	return &RegisterLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *RegisterLogic) Register(in *pb.RegisterReq) (*pb.RegisterResp, error) {
	// 通过 Mobile 查看数据库中是否已存在
	userRet, err := l.svcCtx.UserModel.FindOneByMobile(l.ctx, in.Mobile)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("%w: user(mobile: %s) register: %w", gserr.ErrServerCommon, in.Mobile, err)
	}
	if userRet != nil {
		return nil, fmt.Errorf("%w: user has been resgister: mobile(%s) is exists", gserr.ErrUserExist, in.Mobile)
	}

	// 不存在记录时，插入新记录
	// 对密码进行加密处理
	pwd, err := crypto.EncryptedPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("encrypt password failed: %w", err)
	}

	user := &model.User{
		Username: "",
		Mobile:   in.Mobile,
		Password: pwd,
	}
	ret, err := l.svcCtx.UserModel.Insert(l.ctx, user)
	if err != nil {
		return nil, fmt.Errorf("user(%s) register failed: %v", in.Mobile, err)
	}

	uid, err := ret.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("user(%s) register failed: %v", in.Mobile, err)
	}

	// Token
	generateTokenLogic := NewGenerateTokenLogic(l.ctx, l.svcCtx)
	tokenResp, err := generateTokenLogic.GenerateToken(&pb.GenerateTokenReq{UserId: uid})
	if err != nil {
		return nil, fmt.Errorf("generate token failed: %w", err)
	}

	return &pb.RegisterResp{
		UserId:       uid,
		Token:        tokenResp.Token,
		ExpireAfter:  tokenResp.ExpireAfter,
		RefreshAfter: tokenResp.RefreshAfter,
	}, nil
}
