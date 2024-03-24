package gfcache

import (
	"context"
	"database/sql"
	"github.com/jmoiron/sqlx"
	"reflect"
)

type (
	TransactFunc func(tx *sqlx.Tx) error
)

type SqlxDB struct {
	db *sqlx.DB
}

func NewSqlxDB(db *sqlx.DB) *SqlxDB {
	return &SqlxDB{db: db}
}

func (d *SqlxDB) Query(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	kind := reflect.Indirect(reflect.ValueOf(dest)).Kind()

	var err error
	switch kind {
	case reflect.Slice:
		err = d.db.SelectContext(ctx, dest, query, args...)
	default:
		err = d.db.GetContext(ctx, dest, query, args...)
	}
	if err != nil {
		return err
	}

	return nil
}

func (d *SqlxDB) Execute(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return d.db.ExecContext(ctx, query, args...)
}

func (d *SqlxDB) Transact(ctx context.Context, txFns ...TransactFunc) (retErr error) {
	tx, err := d.db.BeginTxx(ctx, nil)
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
