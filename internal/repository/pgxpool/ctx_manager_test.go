package pgxpool

import (
	"context"
	"testing"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/Ranik23/url-shortener/internal/repository/mock"
	"github.com/stretchr/testify/assert"
)

func TestPgxCtxManager_ByKey(t *testing.T) {
	ctx := context.Background()
	manager := NewPgxCtxManager(nil)

	key := repository.CtxKey{}

	tr := manager.ByKey(ctx, key)
	assert.Nil(t, tr)

	mockTx := mock.NewTransaction(t)

	ctx = context.WithValue(ctx, key, mockTx)
	tr = manager.ByKey(ctx, key)
	assert.Equal(t, mockTx, tr)
}
