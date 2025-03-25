package postgres

import (
	"context"
	"errors"
	"log/slog"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)


type txKey struct{}

type postgresTxManager struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewTxManager(pool *pgxpool.Pool, log *slog.Logger) repository.TxManager {
	return &postgresTxManager{pool: pool, logger: log}
}

func (t *postgresTxManager) GetExecutor(ctx context.Context) repository.Executor {
	if tx, ok := ctx.Value(txKey{}).(repository.Executor); ok {
		return tx
	}
	return t.pool
}

func (t *postgresTxManager) WithTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode, fn func(ctx context.Context) error) error {
	opts := pgx.TxOptions{
		IsoLevel:   isoLevel,
		AccessMode: accessMode,
	}

	var err error

	tx, err := t.pool.BeginTx(ctx, opts)
	if err != nil {
		t.logger.Error("Failed to begin transaction", slog.String("error", err.Error()))
		return err
	}

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil && !errors.Is(rollbackErr, pgx.ErrTxClosed) {
				t.logger.Error("Failed to rollback transaction", slog.String("error", err.Error()))
			}
		}

	}()

	ctx = context.WithValue(ctx, txKey{}, tx)
	if err = fn(ctx); err != nil {
		return err
	}

	err = tx.Commit(ctx)
	if err != nil {
		t.logger.Error("Failed to commit transaction", slog.String("error", err.Error()))
	}
	return err
}