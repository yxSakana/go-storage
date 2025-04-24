package model

import (
	"context"

	"github.com/zeromicro/go-zero/core/stores/cache"
	"github.com/zeromicro/go-zero/core/stores/sqlx"
)

var _ UserFileRelModel = (*customUserFileRelModel)(nil)

type (
	userFileRelTransHandle func(context context.Context, session sqlx.Session) error

	// UserFileRelModel is an interface to be customized, add more methods here,
	// and implement the added methods in customUserFileRelModel.
	UserFileRelModel interface {
		userFileRelModel

		Trans(ctx context.Context, fn userFileRelTransHandle) error
	}

	customUserFileRelModel struct {
		*defaultUserFileRelModel
	}
)

// NewUserFileRelModel returns a model for the database table.
func NewUserFileRelModel(conn sqlx.SqlConn, c cache.CacheConf, opts ...cache.Option) UserFileRelModel {
	return &customUserFileRelModel{
		defaultUserFileRelModel: newUserFileRelModel(conn, c, opts...),
	}
}

func (m *customUserFileRelModel) Trans(ctx context.Context, fn userFileRelTransHandle) error {
	return m.TransactCtx(ctx, func(context context.Context, session sqlx.Session) error {
		return fn(context, session)
	})
}
