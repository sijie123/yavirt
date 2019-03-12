package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/projecteru2/yavirt/db/mysql"
	"github.com/projecteru2/yavirt/db/types"
)

var pool types.Pool = NewSimplePool(8, mysql.Connect)

func Insert(obj interface{}, table string, fields ...string) (sql.Result, error) {
	return pool.Insert(obj, table, fields...)
}

func Get(obj interface{}, query string, args ...interface{}) error {
	return pool.Get(obj, query, args...)
}

func Select(obj interface{}, query string, args ...interface{}) error {
	return pool.Select(obj, query, args...)
}

func Exec(query string, args ...interface{}) (sql.Result, error) {
	return pool.Exec(query, args...)
}

func BeginTx(ctx context.Context) (*sqlx.Tx, func(error), error) {
	return pool.BeginTx(ctx)
}

func SelectForUpdate(ctx context.Context, obj interface{}, selectQuery string, getUpdateQuery func(selectedObj interface{}) string) error {
	return pool.SelectForUpdate(ctx, obj, selectQuery, getUpdateQuery)
}
