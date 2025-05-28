package fileUpoad

import (
	"context"
	"encoding/json"
	"errors"

	"github.com/zeromicro/go-zero/core/logx"

	"go-storage/app/file/cmd/api/internal/svc"
)

type MergeConsumer struct {
	logx.Logger
	ctx    context.Context
	svcCtx *svc.ServiceContext
}

func NewMergeConsumer(ctx context.Context, svcCtx *svc.ServiceContext) *MergeConsumer {
	return &MergeConsumer{
		Logger: logx.WithContext(ctx),
		ctx:    ctx,
		svcCtx: svcCtx,
	}
}

func (l *MergeConsumer) Consume(ctx context.Context, _, val string) error {
	var mess MergeMessage
	if err := json.Unmarshal([]byte(val), &mess); err != nil {
		l.Logger.Errorf("Merge file Consumer: %+v", err)
		return err
	}
	if mess.FileHash == "" || mess.UserId == 0 {
		err := errors.New("merge file Consumer: failed to unmarshal message")
		return err
	}
	meta, err := l.svcCtx.UploadManager.GetMeta(l.ctx, mess.FileHash)
	if err != nil || meta.FileHash == "" {
		if err == nil {
			err = errors.New("merge file Consumer: can't to read file mete into redis")
		}
		l.Logger.Errorf("Merge file Consumer: %+v", err)
		return err
	}

	err = l.svcCtx.MergePool.Submit(func() {
		cl := NewCompleteUploadLogic(ctx, l.svcCtx)
		err = cl.MergeChunks(meta, mess.UserId)
		if err != nil {
			l.Logger.Errorf("Merge file Consumer: %+v", err)
		}
	})
	if err != nil {
		l.Logger.Errorf("Merge file Consumer: %+v", err)
	}
	return err
}
