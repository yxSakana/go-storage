package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ FileMetaModel = (*customFileMetaModel)(nil)

type (
	fileMetaTransHandle func(context context.Context, session sqlx.Session) error

	// FileMetaModel is an interface to be customized, add more methods here,
	// and implement the added methods in customFileMetaModel.
	FileMetaModel interface {
		fileMetaModel

		Trans(ctx context.Context, fn fileMetaTransHandle) error
	}

	customFileMetaModel struct {
		*defaultFileMetaModel
	}
)

// NewFileMetaModel returns a model for the database table.
func NewFileMetaModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) FileMetaModel {
	return &customFileMetaModel{
		defaultFileMetaModel: newFileMetaModel(conn, c, opts...),
	}
}

func (m *customFileMetaModel) Trans(ctx context.Context, fn fileMetaTransHandle) error {
	return m.TransactCtx(ctx, func(context context.Context, session sqlx.Session) error {
		return fn(context, session)
	})
}
