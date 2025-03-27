package pgxpool

import (
	"time"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

type pgxSettings struct {
	txOpts *pgx.TxOptions
	ctxKey struct{}
}

func (p pgxSettings) CtxKey() repository.CtxKey {
	return p.ctxKey
}

func (p *pgxSettings) EnrichBy(external repository.Settings) repository.Settings {
	if ext, ok := external.(*pgxSettings); ok {
		p.txOpts = ext.txOpts
	}
	return p
}

func (p pgxSettings) TimeOutOrNil() *time.Duration {
	panic("unimplemented")
}

func NewPgxSettings() repository.Settings {
	txOpts := &pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	pgxSettings := &pgxSettings{
		txOpts: txOpts,
		ctxKey: struct{}{},
	}

	return pgxSettings
}
