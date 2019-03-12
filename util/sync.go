package util

import "sync/atomic"

type Once struct {
	done int32
}

func (o *Once) Do(fn func() error) (err error) {
	if !atomic.CompareAndSwapInt32(&o.done, 0, 1) {
		return
	}

	if err = fn(); err != nil {
		// rollback
		o.done = 0
	}

	return
}

type AtomicInt64 struct {
	i64 int64
}

func (a *AtomicInt64) Int64() int64 {
	return atomic.LoadInt64(&a.i64)
}

func (a *AtomicInt64) Incr() int64 {
	return a.Add(1)
}

func (a *AtomicInt64) Add(v int64) int64 {
	return atomic.AddInt64(&a.i64, v)
}

func (a *AtomicInt64) Set(v int64) {
	atomic.StoreInt64(&a.i64, v)
}
