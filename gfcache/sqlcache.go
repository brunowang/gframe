package gfcache

import (
	"context"
	"database/sql"
)

type (
	CacheClearFunc func(cache Cache) error
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
