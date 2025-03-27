package inmemory

import (
	"context"
	"github.com/Ranik23/url-shortener/internal/repository"
)

type InMemoryTxManager struct {}

func (i *InMemoryTxManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func (i *InMemoryTxManager) DoWithSettings(ctx context.Context, settings repository.Settings, fn func(context.Context) error) error {
	return fn(ctx)
}

func NewInMemoryTxManager() repository.TxManager {
	return &InMemoryTxManager{}
}
