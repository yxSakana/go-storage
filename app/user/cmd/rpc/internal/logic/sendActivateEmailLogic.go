package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/email/cmd/rpc/email"
	"go-storage/app/user/cmd/rpc/internal/svc"
	"go-storage/app/user/cmd/rpc/pb"
)

type SendActivateEmailLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendActivateEmailLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendActivateEmailLogic {
	return &SendActivateEmailLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendActivateEmailLogic) SendActivateEmail(in *pb.SendActivateEmailReq) (*pb.SendActivateEmailResp, error) {
	if !l.svcCtx.RegisterCache.Exist(in.Email) {
		return nil, fmt.Errorf("email is not exist: %s", in.Email)
	}

	info, err := l.svcCtx.RegisterCache.Load(l.ctx, in.Email)
	if err != nil {
		return nil, err
	}

	exp := time.Duration(svc.CacheRegisterInfoExpireSec) * time.Second
	err = l.Send(*info, exp, l.svcCtx.Config.JwtAuth.AccessSecret)

	return &pb.SendActivateEmailResp{}, err
}

func (l *SendActivateEmailLogic) Send(info svc.RegisterCacheInfo, exp time.Duration, secret string) error {
	verifyToken, err := GenerateVerifyToken(info, exp, secret)
	if err != nil {
		return err
	}
	verifyUrl := fmt.Sprintf("%s/api/v1/user/activate?token=%s", "127.0.0.1:5001", verifyToken)

	l.Logger.Infof("verify url:", verifyUrl)
	go func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		_, err := l.svcCtx.EmailRpc.Send(ctx, &email.SendReq{
			To:      info.Email,
			Subject: "Go-Storage Verify account",
			Body:    fmt.Sprintf("<a href=%s>Click to verify account</a>", verifyUrl),
		})
		if err != nil {
			l.Logger.Error(err)
		}
	}()
	return nil
}
