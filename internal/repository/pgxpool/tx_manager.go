package pgxpool

import (
	"context"
	"log/slog"

	"github.com/Ranik23/url-shortener/internal/repository"
	pool "github.com/jackc/pgx/v5/pgxpool"
)

type txManager struct {
	pool   		*pool.Pool
	pgxSettings repository.Settings
	logger 		*slog.Logger
}

func NewTxManager(pool *pool.Pool, log *slog.Logger, pgxSettings repository.Settings) repository.TxManager {
	return &txManager{
		pool: pool,
		logger: log,
		pgxSettings: pgxSettings,
	}
}

func (p *txManager) Do(ctx context.Context, fn func(context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	pgxTx := NewTransaction(tx)

	newCtx := context.WithValue(ctx, p.pgxSettings.CtxKey(), pgxTx)


	if err := fn(newCtx); err != nil {
		pgxTx.Rollback(ctx)
		return err
	}

	return pgxTx.Commit(ctx)
}


func (p *txManager) DoWithSettings(ctx context.Context, pgxsettings repository.Settings, fn func(context.Context) error) error {
	tx, err := p.pool.Begin(ctx)
	if err != nil {
		return err
	}

	pgxTx := NewTransaction(tx)

	newCtx := context.WithValue(ctx, p.pgxSettings.CtxKey(), pgxTx)


	if err := fn(newCtx); err != nil {
		pgxTx.Rollback(ctx)
		return err
	}

	return pgxTx.Commit(ctx)
}