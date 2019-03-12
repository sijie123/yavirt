package types

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

type Conn interface {
	Operation

	Expired() bool
	Close() error
}

type Operation interface {
	Insert(obj interface{}, table string, fields ...string) (sql.Result, error)
	Get(obj interface{}, query string, args ...interface{}) error
	Exec(query string, args ...interface{}) (sql.Result, error)
	Select(obj interface{}, query string, args ...interface{}) error
	BeginTx(context.Context) (tx *sqlx.Tx, closeTx func(error), err error)
	SelectForUpdate(ctx context.Context, obj interface{}, selectQuery string, getUpdateQuery func(selectedObj interface{}) string) error
}

type Pool Operation
