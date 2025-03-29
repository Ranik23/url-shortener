package inmemory

import (
	"context"
	"github.com/Ranik23/url-shortener/internal/repository"
)

type txManager struct {}

func (i *txManager) Do(ctx context.Context, fn func(context.Context) error) error {
	return fn(ctx)
}

func (i *txManager) DoWithSettings(ctx context.Context, settings repository.Settings, fn func(context.Context) error) error {
	return fn(ctx)
}

func NewTxManager() repository.TxManager {
	return &txManager{}
}
