package inmemory

import (
	"context"
	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

type InMemoryTxManager struct {
	repo *inMemoryRepository
}

func (i *InMemoryTxManager) GetExecutor(ctx context.Context) repository.Executor {
	return repo
}

func (i *InMemoryTxManager) WithTx(ctx context.Context, isoLevel pgx.TxIsoLevel, accessMode pgx.TxAccessMode, fn func(ctx context.Context) error) error {
	return fn(ctx)
}

func NewInMemoryTxManager() repository.TxManager {
	return &InMemoryTxManager{}
}
