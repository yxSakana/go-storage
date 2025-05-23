// Code generated by goctl. DO NOT EDIT.
// goctl 1.8.2
// Source: user.proto

package user

import (
	"context"

	"go-storage/app/user/cmd/rpc/pb"

	"github.com/zeromicro/go-zero/zrpc"
	"google.golang.org/grpc"
)

type (
	GenerateTokenReq  = pb.GenerateTokenReq
	GenerateTokenResp = pb.GenerateTokenResp
	GetUserInfoReq    = pb.GetUserInfoReq
	GetUserInfoResp   = pb.GetUserInfoResp
	LoginReq          = pb.LoginReq
	LoginResp         = pb.LoginResp
	RegisterReq       = pb.RegisterReq
	RegisterResp      = pb.RegisterResp
	User              = pb.User

	UserZrpcClient interface {
		Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error)
		Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error)
		GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error)
		GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenResp, error)
	}

	defaultUserZrpcClient struct {
		cli zrpc.Client
	}
)

func NewUserZrpcClient(cli zrpc.Client) UserZrpcClient {
	return &defaultUserZrpcClient{
		cli: cli,
	}
}

func (m *defaultUserZrpcClient) Register(ctx context.Context, in *RegisterReq, opts ...grpc.CallOption) (*RegisterResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.Register(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) Login(ctx context.Context, in *LoginReq, opts ...grpc.CallOption) (*LoginResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.Login(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) GetUserInfo(ctx context.Context, in *GetUserInfoReq, opts ...grpc.CallOption) (*GetUserInfoResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GetUserInfo(ctx, in, opts...)
}

func (m *defaultUserZrpcClient) GenerateToken(ctx context.Context, in *GenerateTokenReq, opts ...grpc.CallOption) (*GenerateTokenResp, error) {
	client := pb.NewUserClient(m.cli.Conn())
	return client.GenerateToken(ctx, in, opts...)
}
