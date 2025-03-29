package pgxpool

import (
	"context"
	"sync"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers"
	pg"github.com/jackc/pgx/v5"
)

type transaction struct {
	mu 			sync.Mutex
	tx 			pg.Tx
	isClosed 	*drivers.IsClosed
}

func (t *transaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *transaction) IsActive() bool {
	return t.isClosed.IsActive()
}

func (t *transaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *transaction) Transaction() interface{} {
	return t.tx
}

func NewTransaction(tx pg.Tx) repository.Transaction {
	return &transaction{
		isClosed: drivers.NewIsClosed(),
		tx: tx,
	}
}
