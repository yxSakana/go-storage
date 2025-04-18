package logic

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/token"
	"go-storage/app/user/cmd/rpc/internal/svc"
	"go-storage/app/user/cmd/rpc/pb"
)

type GenerateTokenLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewGenerateTokenLogic(ctx context.Context, svcCtx *svc.ServiceContext) *GenerateTokenLogic {
	return &GenerateTokenLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *GenerateTokenLogic) GenerateToken(in *pb.GenerateTokenReq) (*pb.GenerateTokenResp, error) {
	now := time.Now()
	expire := l.svcCtx.Config.JwtAuth.AccessExpire
	secret := l.svcCtx.Config.JwtAuth.AccessSecret
	claims := token.GsClaims{
		UserId: in.UserId,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(now.Add(time.Duration(expire) * time.Second)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := jwtToken.SignedString([]byte(secret))
	if err != nil {
		return nil, fmt.Errorf("jwt token sign err: %w", err)
	}

	return &pb.GenerateTokenResp{
		Token:        tokenString,
		ExpireAfter:  expire,
		RefreshAfter: expire / 2,
	}, nil
}
