package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"

	"github.com/projecteru2/yavirt/db/types"
	"github.com/projecteru2/yavirt/test/mock"
)

type MockPool struct {
	mock.Mock
}

func NewMockPool() (types.Pool, func()) {
	var p = &MockPool{}
	return p, mockPool(p)
}

func mockPool(p types.Pool) func() {
	var old = pool

	pool = p

	return func() {
		pool = old
	}
}

func (p *MockPool) Insert(obj interface{}, table string, fields ...string) (sql.Result, error) {
	return &Result{}, nil
}

func (p *MockPool) Get(obj interface{}, query string, args ...interface{}) error {
	return nil
}

func (p *MockPool) Exec(query string, args ...interface{}) (sql.Result, error) {
	return &Result{}, nil
}

func (p *MockPool) Select(obj interface{}, query string, args ...interface{}) error {
	return nil
}

func (p *MockPool) BeginTx(ctx context.Context) (*sqlx.Tx, func(error), error) {
	return &sqlx.Tx{}, func(error) {}, nil
}

func (p *MockPool) SelectForUpdate(ctx context.Context, obj interface{}, selectQuery string, getUpdateQuery func(selectedObj interface{}) string) error {
	return nil
}

type Result struct{}

func (r *Result) LastInsertId() (int64, error) { return 999, nil }
func (r *Result) RowsAffected() (int64, error) { return 1, nil }
