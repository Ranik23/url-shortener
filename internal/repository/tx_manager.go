package repository

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)


type Executor interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
}

type TxManager interface {
	GetExecutor(ctx context.Context) Executor
	WithTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode, fn func(ctx context.Context) error) error
}