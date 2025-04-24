package response

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/zeromicro/go-zero/core/logx"
	"github.com/zeromicro/go-zero/rest/httpx"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"go-storage/pkg/gserr"
)

type Response struct {
	Code uint32      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func HttpResult(ctx context.Context, w http.ResponseWriter, resp interface{}, err error) {
	var (
		code uint32 = 0
		msg  string = "success"
	)
	if err == nil {
		httpx.WriteJson(w, http.StatusOK, Response{
			Code: code,
			Msg:  msg,
			Data: resp,
		})
		return
	} else {
		code = gserr.ServerCommonError
		msg = err.Error()
	}

	// 自定义错误
	errToHttp := gserr.ErrUnknown
	if errors.As(err, &errToHttp) {
		code = errToHttp.Code()
		msg = errToHttp.Message()
		if errors.Is(errToHttp, gserr.ErrAttachedMsgError) {
			msg = err.Error()
		}
	}

	// RPC错误
	if st, ok := statusFromError(err); ok {
		code = uint32(st.Code())
		msg = st.Message()
		// RPC中的自定义code
		//if m, ok := gserr.MsgFromCode(code); ok {
		//	msg = m
		//}
	}
	httpx.WriteJson(w, http.StatusInternalServerError, Response{
		Code: code,
		Msg:  msg,
		Data: nil,
	})
	logx.WithContext(ctx).Errorf("[API-ERROR]: %v", err)
}

func ParamError(ctx context.Context, w http.ResponseWriter, err error) {
	logx.WithContext(ctx).Errorf("[API-ERROR]: %v", err)
	errToHttp := gserr.ErrRequestParam
	httpx.WriteJson(w, http.StatusBadRequest, Response{
		Code: errToHttp.Code(),
		Msg:  fmt.Sprintf("%s: %v", errToHttp.Message(), err),
	})
}

func UnauthorizedCallback(w http.ResponseWriter, _ *http.Request, err error) {
	httpx.WriteJson(w, http.StatusUnauthorized, Response{
		Code: http.StatusUnauthorized,
		Msg:  fmt.Sprintf("未登录: %s", err.Error()),
		Data: nil,
	})
}

func statusFromError(err error) (s *status.Status, ok bool) {
	if err == nil {
		return nil, true
	}
	type grpcstatus interface{ GRPCStatus() *status.Status }
	if gs, ok := err.(grpcstatus); ok {
		grpcStatus := gs.GRPCStatus()
		if grpcStatus == nil {
			return status.New(codes.Unknown, err.Error()), false
		}
		return grpcStatus, true
	}
	var gs grpcstatus
	if errors.As(err, &gs) {
		grpcStatus := gs.GRPCStatus()
		if grpcStatus == nil {
			return status.New(codes.Unknown, err.Error()), false
		}
		p := grpcStatus.Proto()
		p.Message = grpcStatus.Message()
		return status.FromProto(p), true
	}
	return status.New(codes.Unknown, err.Error()), false
}
