package interceptor

import (
	"context"
	"errors"
	"github.com/zeromicro/go-zero/core/logx"
	"go-storage/pkg/gserr"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func RpcErrCovAndLoggerInterceptor(ctx context.Context, req any,
	_ *grpc.UnaryServerInfo, handler grpc.UnaryHandler,
) (resp any, err error) {

	resp, err = handler(ctx, req)

	if err != nil {
		logx.WithContext(ctx).Errorf("[gRPC-ERROR]: %+v", err)

		var customErr *gserr.Error
		if errors.As(err, &customErr) {
			err = status.Error(codes.Code(customErr.Code()), customErr.Error())
		}
	}

	return resp, err
}
