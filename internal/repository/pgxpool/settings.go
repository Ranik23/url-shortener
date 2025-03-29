package pgxpool

import (
	"time"

	"github.com/Ranik23/url-shortener/internal/repository"
	"github.com/jackc/pgx/v5"
)

type settings struct {
	txOpts *pgx.TxOptions
	ctxKey struct{}
}

func (p settings) CtxKey() repository.CtxKey {
	return p.ctxKey
}

func (p *settings) EnrichBy(external repository.Settings) repository.Settings {
	if ext, ok := external.(*settings); ok {
		p.txOpts = ext.txOpts
	}
	return p
}

func (p settings) TimeOutOrNil() *time.Duration {
	panic("unimplemented")
}

func NewSettings() repository.Settings {
	txOpts := &pgx.TxOptions{
		IsoLevel:   pgx.Serializable,
		AccessMode: pgx.ReadWrite,
	}

	settings := &settings{
		txOpts: txOpts,
		ctxKey: struct{}{},
	}

	return settings
}
