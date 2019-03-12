package util

import (
	"sync"
	"testing"

	"github.com/projecteru2/yavirt/test/assert"
)

func TestOnce(t *testing.T) {
	var once Once
	var wg sync.WaitGroup
	var n int

	for i := 0; i < 64; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NilErr(t, once.Do(func() error {
				n++
				return nil
			}))
		}()
	}
	wg.Wait()

	assert.Equal(t, 1, n)
}

func TestAtomicInt64(t *testing.T) {
	var wg sync.WaitGroup
	var atom AtomicInt64
	var n = 64

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			atom.Incr()
		}()
	}
	wg.Wait()

	assert.Equal(t, int64(n), atom.Int64())
}
