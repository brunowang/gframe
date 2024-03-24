package gfcache

import (
	"context"
	"database/sql"
	"errors"
)

var (
	RecordNotFound = errors.New("record not found")
	CacheNotFound  = errors.New("cache not found")
)

type Cache interface {
	HasCache(key string) bool
	GetCache(key string, data interface{}) error
	SetCache(key string, data interface{}) error
	DelCache(key string) error
}

type SourceDB interface {
	Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error
	Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type SourceDBTx[T Transaction] interface {
	Transact(ctx context.Context, txFn TransactFunc[T]) error
	SourceDB() SourceDB
}

type Transaction interface {
	Commit() error
	Rollback() error
}
