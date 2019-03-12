package volume

import (
	"fmt"
	"path/filepath"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/sh"
	"github.com/projecteru2/yavirt/util"
	"github.com/projecteru2/yavirt/virt/common"
)

type Bot interface {
	Close() error
	Undefine() error
	CopyFrom(src common.ResourceDomain) error
}

type virtVol struct {
	vol   *Volume
	flock *util.Flock
}

func newVirtVol(vol *Volume) (Bot, error) {
	var virt = &virtVol{vol: vol}
	virt.flock = virt.newFlock()

	if err := virt.flock.Trylock(); err != nil {
		return nil, errors.Trace(err)
	}

	return virt, nil
}

func (v *virtVol) Undefine() error {
	return sh.Remove(v.vol.Filepath())
}

func (v *virtVol) Close() error {
	v.flock.Close()
	return nil
}

func (v *virtVol) newFlock() *util.Flock {
	var fn = fmt.Sprintf("%s.flock", v.vol.Name())
	var fpth = filepath.Join(config.Conf.VirtFlockDir, fn)
	return util.NewFlock(fpth)
}

func (v *virtVol) CopyFrom(src common.ResourceDomain) error {
	if err := sh.Copy(src.Filepath(), v.vol.Filepath()); err != nil {
		return errors.Trace(err)
	}
	return nil
}
