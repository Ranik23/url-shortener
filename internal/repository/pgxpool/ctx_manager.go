package pgxpool

import (
	"context"

	"github.com/Ranik23/url-shortener/internal/repository"
	pool "github.com/jackc/pgx/v5/pgxpool"
)

type ctxManager struct {
	pool *pool.Pool
}

func (p *ctxManager) ByKey(ctx context.Context, key repository.CtxKey) repository.Transaction {
	tx, ok := ctx.Value(key).(repository.Transaction)
	if !ok {
		return nil
	}
	return tx
}


func (p *ctxManager) Default(ctx context.Context) repository.Transaction {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return nil
	}
	return NewTransaction(tx)
}

func NewCtxManager(pool *pool.Pool) repository.CtxManager {
	return &ctxManager{
		pool: pool,
	}
}
