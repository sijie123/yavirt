package db

import (
	"container/list"
	"context"
	"database/sql"
	"sync"

	"github.com/jmoiron/sqlx"
	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/db/types"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/metric"
)

type SimplePool struct {
	sync.Mutex
	conns   list.List
	connect Connect
	cap     int
	created int
	idle    int
}

type Connect func() (types.Conn, error)

func NewSimplePool(cap int, connect Connect) *SimplePool {
	return &SimplePool{
		connect: connect,
		cap:     cap,
	}
}

func (p *SimplePool) Insert(obj interface{}, table string, fields ...string) (sql.Result, error) {
	conn, put, err := p.getConn()
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer put()

	res, err := conn.Insert(obj, table, fields...)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return res, nil
}

func (p *SimplePool) Get(obj interface{}, query string, args ...interface{}) error {
	var conn, put, err = p.getConn()
	if err != nil {
		return errors.Trace(err)
	}

	defer put()

	return conn.Get(obj, query, args...)
}

func (p *SimplePool) Select(obj interface{}, query string, args ...interface{}) error {
	var conn, put, err = p.getConn()
	if err != nil {
		return errors.Trace(err)
	}

	defer put()

	return conn.Select(obj, query, args...)
}

func (p *SimplePool) Exec(query string, args ...interface{}) (sql.Result, error) {
	var conn, put, err = p.getConn()
	if err != nil {
		return nil, errors.Trace(err)
	}

	defer put()

	return conn.Exec(query, args...)
}

func (p *SimplePool) BeginTx(ctx context.Context) (*sqlx.Tx, func(error), error) {
	conn, put, err := p.getConn()
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	tx, closeTx, err := conn.BeginTx(ctx)
	if err != nil {
		return nil, nil, errors.Trace(err)
	}

	var clean = func(err error) {
		closeTx(err)
		put()
	}

	return tx, clean, nil
}

func (p *SimplePool) SelectForUpdate(ctx context.Context,
	obj interface{},
	selectQuery string,
	getUpdateQuery func(interface{}) string) error {

	var conn, put, err = p.getConn()
	if err != nil {
		return errors.Trace(err)
	}

	defer put()

	return conn.SelectForUpdate(ctx, obj, selectQuery, getUpdateQuery)
}

func (p *SimplePool) getConn() (conn types.Conn, put func(), err error) {
	if conn, err = p.GetConn(); err != nil {
		return
	}

	put = func() {
		if e := p.PutConn(conn); e != nil {
			log.Warnf(errors.ErrorStack(e))
			metric.IncrError()
		}
	}

	return
}

func (p *SimplePool) GetConn() (types.Conn, error) {
	p.Lock()
	defer p.Unlock()

	if p.idle < 1 {
		var conn, err = p.connect()
		if err != nil {
			return nil, errors.Trace(err)
		}

		p.created++

		return conn, nil
	}

	var elem = p.conns.Back()
	if elem == nil {
		return nil, errors.Errorf("empty pool")
	}

	defer func() {
		p.idle--
	}()

	var conn, ok = p.conns.Remove(elem).(types.Conn)
	if !ok {
		return nil, errors.Errorf("invalid types.Conn")
	}

	return conn, nil
}

func (p *SimplePool) PutConn(conn types.Conn) error {
	p.Lock()
	defer p.Unlock()

	if p.idle >= p.cap {
		return errors.Errorf("pool is full")
	}

	if conn.Expired() {
		return conn.Close()
	}

	p.idle++

	p.conns.PushFront(conn)

	return nil
}
