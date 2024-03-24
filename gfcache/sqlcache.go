package gfcache

import (
	"context"
	"database/sql"
)

type (
	CacheClearFunc              func(cache Cache) error
	TransactFunc[T Transaction] func(tx T) error
)

type SqlCache struct {
	db    SourceDB
	cache Cache
}

func NewSqlCache(db SourceDB, cache Cache) *SqlCache {
	return &SqlCache{db: db, cache: cache}
}

func (c *SqlCache) Query(ctx context.Context, cacheKey string, dest interface{}, query string, args ...interface{}) error {
	if err := c.cache.GetCache(cacheKey, dest); err == nil {
		return nil
	}
	if err := c.db.Query(ctx, dest, query, args...); err != nil {
		return err
	}
	_ = c.cache.SetCache(cacheKey, dest)
	return nil
}

func (c *SqlCache) Execute(ctx context.Context, clearFn CacheClearFunc, query string, args ...interface{}) (sql.Result, error) {
	res, err := c.db.Execute(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if err := clearFn(c.cache); err != nil {
		return nil, err
	}
	return res, nil
}

type SqlCacheTx[T Transaction] struct {
	SqlCache
	db SourceDBTx[T]
}

func NewSqlCacheTx[T Transaction](db SourceDBTx[T], cache Cache) *SqlCacheTx[T] {
	return &SqlCacheTx[T]{db: db, SqlCache: SqlCache{db: db.SourceDB(), cache: cache}}
}

func (c *SqlCacheTx[T]) Transact(ctx context.Context, clearFn CacheClearFunc, txFn TransactFunc[T]) error {
	if err := c.db.Transact(ctx, func(tx T) error {
		return txFn(tx)
	}); err != nil {
		return err
	}
	if err := clearFn(c.cache); err != nil {
		return err
	}
	return nil
}
