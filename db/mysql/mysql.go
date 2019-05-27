package mysql

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/db/types"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/metric"
)

const Driver = "mysql"

type Mysql struct {
	dbx *sqlx.DB
	gen *generator
	dsn *Dsn

	birth time.Time
	life  time.Duration
}

func Connect() (types.Conn, error) {
	var conn = &Mysql{
		gen: newGenerator(),
		dsn: &Dsn{
			User:     config.Conf.MysqlUser,
			Password: config.Conf.MysqlPassword,
			Addr:     config.Conf.MysqlAddr,
			DB:       config.Conf.MysqlDB,
		},
		birth: time.Now(),
		life:  time.Hour * 2,
	}

	var ctx, cancel = conn.newContext()
	defer cancel()

	var dbx, err = sqlx.ConnectContext(ctx, Driver, conn.dsn.String())
	if err != nil {
		return nil, errors.Trace(err)
	}

	conn.dbx = dbx
	conn.birth = time.Now()

	return conn, nil
}

func (m *Mysql) Get(obj interface{}, query string, args ...interface{}) error {
	var ctx, cancel = m.newContext()
	defer cancel()

	if err := m.dbx.GetContext(ctx, obj, query, args...); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (m *Mysql) Select(obj interface{}, query string, args ...interface{}) error {
	var ctx, cancel = m.newContext()
	defer cancel()

	if err := m.dbx.SelectContext(ctx, obj, query, args...); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (m *Mysql) Insert(obj interface{}, table string, fields ...string) (sql.Result, error) {
	var query = m.gen.genInsertx(table, fields...)
	return m.namedExec(query, obj)
}

func (m *Mysql) Exec(query string, args ...interface{}) (sql.Result, error) {
	var ctx, cancel = m.newContext()
	defer cancel()

	res, err := m.dbx.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return res, nil
}

func (m *Mysql) namedExec(query string, obj interface{}) (sql.Result, error) {
	var ctx, cancel = m.newContext()
	defer cancel()

	res, err := m.dbx.NamedExecContext(ctx, query, obj)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return res, nil
}

func (m *Mysql) SelectForUpdate(ctx context.Context,
	obj interface{},
	selectQuery string,
	getUpdateQuery func(interface{}) string) error {

	var tx, closeTx, err = m.BeginTx(ctx)
	if err != nil {
		return errors.Trace(err)
	}

	defer func() {
		closeTx(err)
	}()

	for {
		if err := tx.GetContext(ctx, obj, selectQuery); err != nil {
			return errors.Annotatef(err, "failed to '%s'", selectQuery)
		}

		var res, err = tx.ExecContext(ctx, getUpdateQuery(obj))
		if err != nil {
			return errors.Trace(err)
		}

		switch affe, err := res.RowsAffected(); {
		case err != nil:
			return errors.Trace(err)
		case affe < 1:
			// The selected rows had been update by another thread.
			continue
		default:
			return nil
		}
	}
}

func (m *Mysql) BeginTx(ctx context.Context) (tx *sqlx.Tx, closeTx func(error), err error) {
	if tx, err = m.dbx.BeginTxx(ctx, nil); err != nil {
		return
	}

	closeTx = func(err error) {
		var xe error
		if err != nil {
			xe = tx.Rollback()
		} else {
			xe = tx.Commit()
		}

		if xe != nil {
			log.Errorf(errors.ErrorStack(xe))
			metric.IncrError()
		}
	}

	return
}

func (m *Mysql) newContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), m.timeout())
}

func (m *Mysql) timeout() time.Duration {
	return config.Conf.MysqlTimeout.Duration()
}

func (m *Mysql) Expired() bool {
	return time.Since(m.birth) >= m.life
}

func (m *Mysql) Close() error {
	return m.dbx.Close()
}
