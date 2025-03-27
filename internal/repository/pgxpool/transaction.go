package pgxpool

import (
	"context"
	"sync"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/avito-tech/go-transaction-manager/trm/v2/drivers"
	pg"github.com/jackc/pgx/v5"
)

type pgxTransaction struct {
	mu 			sync.Mutex
	tx 			pg.Tx
	isClosed 	*drivers.IsClosed
}

func (t *pgxTransaction) Commit(ctx context.Context) error {
	return t.tx.Commit(ctx)
}

func (t *pgxTransaction) IsActive() bool {
	return t.isClosed.IsActive()
}

func (t *pgxTransaction) Rollback(ctx context.Context) error {
	return t.tx.Rollback(ctx)
}

func (t *pgxTransaction) Transaction() interface{} {
	return t.tx
}

func NewTransaction(tx pg.Tx) repository.Transaction {
	return &pgxTransaction{
		isClosed: drivers.NewIsClosed(),
		tx: tx,
	}
}
