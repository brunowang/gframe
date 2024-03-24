package gfcache

import (
	"context"
	"database/sql"
	"errors"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type (
	CacheClearFunc func(cache Cache) error
	TransactFunc   func(tx *sqlx.Tx) error
)

type SqlCache struct {
	db    *sqlx.DB
	cache Cache
}

func NewSqlCache(db *sqlx.DB, cache Cache) *SqlCache {
	return &SqlCache{db: db, cache: cache}
}

func (c *SqlCache) Query(ctx context.Context, cacheKey string, dest interface{}, query string, args ...interface{}) error {
	value := reflect.ValueOf(dest)

	// json.Unmarshal returns errors for these
	if value.Kind() != reflect.Ptr {
		return errors.New("must pass a pointer, not a value, to StructScan destination")
	}
	if value.IsNil() {
		return errors.New("nil pointer passed to StructScan destination")
	}
	value = value.Elem()

	if err := c.cache.GetCache(cacheKey, dest); err == nil {
		return nil
	}

	var err error
	switch value.Kind() {
	case reflect.Slice:
		err = c.db.SelectContext(ctx, dest, query, args...)
	default:
		err = c.db.GetContext(ctx, dest, query, args...)
	}
	if err != nil {
		return err
	}

	_ = c.cache.SetCache(cacheKey, dest)

	return nil
}

func (c *SqlCache) Execute(ctx context.Context, clearFn CacheClearFunc, query string, args ...interface{}) (sql.Result, error) {
	res, err := c.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	if err := clearFn(c.cache); err != nil {
		return nil, err
	}
	return res, nil
}

func (c *SqlCache) Transact(ctx context.Context, txFns ...TransactFunc) (retErr error) {
	tx, err := c.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}
	defer func() {
		if retErr != nil {
			_ = tx.Rollback()
			return
		}
		retErr = tx.Commit()
	}()

	for _, fn := range txFns {
		if err := fn(tx); err != nil {
			return err
		}
	}
	return nil
}
