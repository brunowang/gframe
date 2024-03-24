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
	Transact(ctx context.Context, txFns ...TransactFunc) error
}
