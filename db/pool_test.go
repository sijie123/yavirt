package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/jmoiron/sqlx"

	"github.com/projecteru2/yavirt/db/types"
	"github.com/projecteru2/yavirt/test/assert"
	"github.com/projecteru2/yavirt/test/mock"
)

func TestPool(t *testing.T) {
	var cap = 2
	var pool = NewSimplePool(cap, func() (types.Conn, error) {
		return &mockConn{}, nil
	})

	var n = 3
	var conns = make([]types.Conn, n)
	var err error

	for i := 0; i < n; i++ {
		conns[i], err = pool.GetConn()
		assert.NilErr(t, err)
	}

	assert.Equal(t, n, pool.created)
	assert.Equal(t, 0, pool.idle)

	assert.NilErr(t, pool.PutConn(conns[0]))
	assert.NilErr(t, pool.PutConn(conns[1]))
	assert.Err(t, pool.PutConn(conns[2]))

	assert.Equal(t, n, pool.created)
	assert.Equal(t, cap, pool.idle)
}

func TestPoolExpired(t *testing.T) {
	var cap = 2
	var pool = NewSimplePool(cap, func() (types.Conn, error) {
		return &mockConn{expired: true}, nil
	})

	for i := 0; i < cap; i++ {
		var conn, err = pool.GetConn()
		assert.NilErr(t, err)
		assert.NilErr(t, pool.PutConn(conn))
	}

	assert.Equal(t, cap, pool.created)
	assert.Equal(t, 0, pool.idle)
}

func TestPooled(t *testing.T) {
	var cap = 2
	var pool = NewSimplePool(cap, func() (types.Conn, error) {
		return &mockConn{}, nil
	})

	for i := 0; i < 2; i++ {
		var conn, err = pool.GetConn()
		assert.NilErr(t, err)
		assert.Equal(t, 0, pool.idle)
		assert.NilErr(t, pool.PutConn(conn))
		assert.Equal(t, 1, pool.created)
		assert.Equal(t, 1, pool.idle)
	}
}

type mockConn struct {
	mock.Mock
	expired bool
}

func (c *mockConn) Close() error {
	return nil
}

func (c *mockConn) Expired() bool {
	return c.expired
}

func (c *mockConn) Insert(obj interface{}, table string, fields ...string) (sql.Result, error) {
	return &Result{}, nil
}

func (c *mockConn) Get(obj interface{}, query string, args ...interface{}) error {
	return nil
}

func (c *mockConn) Exec(query string, args ...interface{}) (sql.Result, error) {
	return &Result{}, nil
}

func (c *mockConn) Select(obj interface{}, query string, args ...interface{}) error {
	return nil
}

func (c *mockConn) BeginTx(ctx context.Context) (*sqlx.Tx, func(error), error) {
	return &sqlx.Tx{}, func(error) {}, nil
}

func (c *mockConn) SelectForUpdate(ctx context.Context, obj interface{}, selectQuery string, getUpdateQuery func(selectedObj interface{}) string) error {
	return nil
}
