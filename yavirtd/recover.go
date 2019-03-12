package yavirtd

import (
	"sync"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/metric"
	"github.com/projecteru2/yavirt/util"
	"github.com/projecteru2/yavirt/virt/guest"
)

func (s *Server) disasterRecover() error {
	return s.recoverGuests()
}

func (s *Server) recoverGuests() error {
	var wg sync.WaitGroup
	var errCnt util.AtomicInt64

	var guests, err = guest.ListLocals()
	if err != nil {
		return errors.Trace(err)
	}

	for _, g := range guests {
		wg.Add(1)

		go func(g *guest.Guest) {
			defer wg.Done()

			if err := g.SyncState(); err != nil {
				log.Errorf(errors.ErrorStack(err))
				metric.IncrError()
				errCnt.Incr()
			}
		}(g)
	}

	if errCnt.Int64() > 0 {
		return errors.Errorf("some guests recovery failed")
	}

	return nil
}
