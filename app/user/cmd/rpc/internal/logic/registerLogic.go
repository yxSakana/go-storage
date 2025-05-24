package logic

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
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
	// 确保email不存在
	userRet, err := l.svcCtx.UserModel.FindOneByEmail(l.ctx, in.Email)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("%w: user(Email: %s) register: %w", gserr.ErrServerCommon, in.Email, err)
	}
	if userRet != nil {
		return nil, fmt.Errorf("%w: user has been resgister: Email(%s) is exists", gserr.ErrUserExist, in.Email)
	}
	if l.svcCtx.RegisterCache.Exist(in.Email) {
		return nil, fmt.Errorf("%w: user(Email: %s) is exists", gserr.ErrUserExist, in.Email)
	}
	// 暂存到redis: email, code, password
	// TODO: 确保email格式合法
	pwd, err := crypto.EncryptedPassword(in.Password)
	if err != nil {
		return nil, fmt.Errorf("encrypt password failed: %w", err)
	}

	info := svc.RegisterCacheInfo{
		Email:    in.Email,
		Password: pwd,
		Token:    uuid.NewString(),
	}
	if err := l.svcCtx.RegisterCache.Save(l.ctx, &info); err != nil {
		return nil, fmt.Errorf("register cache failed: %w", err)
	}
	// 发送激活链接
	exp := time.Duration(svc.CacheRegisterInfoExpireSec) * time.Second
	sl := NewSendActivateEmailLogic(l.ctx, l.svcCtx)
	err = sl.Send(info, exp, l.svcCtx.Config.JwtAuth.AccessSecret)

	return &pb.RegisterResp{}, err
}
