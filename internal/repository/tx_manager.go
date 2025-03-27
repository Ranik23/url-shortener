package repository

import (
	"context"
	"time"
)

type CtxKey struct{}


type Settings interface {
	EnrichBy(external Settings) Settings
	CtxKey() CtxKey
	TimeOutOrNil() *time.Duration
}

type TxManager interface {
	Do(context.Context, func(context.Context) error) error
	DoWithSettings(context.Context, Settings, func(context.Context) error) error
}

type Transaction interface {
    Commit(context.Context) error 
    Rollback(context.Context) error 
    Transaction() interface{}
 }

 type CtxManager interface {
    Default(context.Context) Transaction
    ByKey(context.Context, CtxKey) Transaction
 }