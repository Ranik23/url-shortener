package pgxpool

import (
	"context"
	"testing"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/Ranik23/url-shortener/internal/repository/mock"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestPgxCtxManager_ByKey(t *testing.T) {
	ctx := context.Background()
	manager := NewPgxCtxManager(nil)

	key := repository.CtxKey{}

	// Проверяем, что если в контексте нет транзакции, то вернется nil
	tr := manager.ByKey(ctx, key)
	assert.Nil(t, tr)


	ctrl := gomock.NewController(t)

	// Добавляем транзакцию в контекст и проверяем, что она корректно извлекается
	mockTr := mock.NewMockTransaction(ctrl)
	ctx = context.WithValue(ctx, key, mockTr)
	tr = manager.ByKey(ctx, key)
	assert.Equal(t, mockTr, tr)
}
