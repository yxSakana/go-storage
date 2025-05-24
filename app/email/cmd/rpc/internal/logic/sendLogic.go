package logic

import (
	"context"
	"gopkg.in/gomail.v2"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/email/cmd/rpc/internal/svc"
	"go-storage/app/email/cmd/rpc/pb"
)

type SendLogic struct {
	ctx    context.Context
	svcCtx *svc.ServiceContext
	logx.Logger
}

func NewSendLogic(ctx context.Context, svcCtx *svc.ServiceContext) *SendLogic {
	return &SendLogic{
		ctx:    ctx,
		svcCtx: svcCtx,
		Logger: logx.WithContext(ctx),
	}
}

func (l *SendLogic) Send(in *pb.SendReq) (*pb.SendResp, error) {
	ec := l.svcCtx.Config.EmailConfig

	m := gomail.NewMessage()
	m.SetHeaders(map[string][]string{
		"From":    {ec.From},
		"To":      {in.To},
		"Subject": {in.Subject},
	})
	m.SetBody("text/html", in.Body)

	d := gomail.NewDialer(ec.Host, ec.Port, ec.Username, ec.Password)
	if err := d.DialAndSend(m); err != nil {
		return nil, err
	}
	return &pb.SendResp{}, nil
}
