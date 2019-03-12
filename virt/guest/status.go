package guest

import (
	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/virt/common"
)

func (g *Guest) setCreating() error {
	return g.setStatus(common.StatusCreating)
}

func (g *Guest) setStarting() error {
	return g.setStatus(common.StatusStarting)
}

func (g *Guest) setStopping() error {
	return g.setStatus(common.StatusStopping)
}

func (g *Guest) setStopped() error {
	return g.setStatus(common.StatusStopped)
}

func (g *Guest) setDestroying() error {
	return g.setStatus(common.StatusDestroying)
}

func (g *Guest) setDestroyed() error {
	return g.setStatus(common.StatusDestroyed)
}

func (g *Guest) setRunning() error {
	return g.setStatus(common.StatusRunning)
}

func (g *Guest) setResizing() error {
	return g.setStatus(common.StatusResizing)
}

func (g *Guest) setMigrating() error {
	return g.setStatus(common.StatusMigrating)
}

func (g *Guest) setStatus(st string) error {
	for _, fn := range []func(string) error{g.vols.UpdateStatus, g.UpdateStatus} {
		if err := fn(st); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (g *Guest) UpdateStatus(st string) (err error) {
	return g.Resource.UpdateStatus(common.TableGuest, st)
}
