package guest

import (
	"context"
	"fmt"
	"path/filepath"
	"time"

	"github.com/juju/errors"
	"github.com/libvirt/libvirt-go"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/util"
)

type Bot interface {
	Close() error
	Create() error
	Boot() error
	Shutdown() error
	Undefine() error
	Migrate() error
}

type virtGuest struct {
	guest *Guest
	virt  *libvirt.Connect
	dom   *domain
	ga    *Agent
	flock *util.Flock
}

func newVirtGuest(guest *Guest) (Bot, error) {
	var virt, err = libvirt.NewConnect("qemu:///system")
	if err != nil {
		return nil, errors.Trace(err)
	}

	var vm = &virtGuest{
		guest: guest,
		virt:  virt,
	}
	vm.dom = vm.domain()
	vm.flock = vm.newFlock()
	vm.ga = NewAgent(vm.guest.sockfile())

	if err := vm.flock.Trylock(); err != nil {
		return nil, errors.Trace(err)
	}

	return vm, nil
}

func (v *virtGuest) Close() (err error) {
	v.flock.Close()

	if err = v.dom.close(); err != nil {
		log.Warnf(errors.ErrorStack(err))
	}

	if _, err = v.virt.Close(); err != nil {
		log.Warnf(errors.ErrorStack(err))
	}

	if err = v.ga.Close(); err != nil {
		log.Warnf(errors.ErrorStack(err))
	}

	return
}

func (v *virtGuest) Migrate() error {
	// TODO
	return nil
}

func (v *virtGuest) Boot() error {
	return util.Invoke([]func() error{
		v.dom.boot,
		v.postBoot,
	})
}

func (v *virtGuest) postBoot() error {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute*30)
	defer cancel()

	if err := v.ga.Ping(ctx); err != nil {
		return errors.Trace(err)
	}

	return v.setNics()
}

func (v *virtGuest) Shutdown() error {
	return v.dom.shutdown()
}

func (v *virtGuest) Undefine() error {
	return util.Invoke([]func() error{
		v.dom.undefine,
		v.guest.sysVol.Undefine,
	})
}

func (v *virtGuest) Create() error {
	return util.Invoke([]func() error{
		v.allocVol,
		v.allocGuest,
	})
}

func (v *virtGuest) allocVol() error {
	if err := v.guest.Image.Cache(); err != nil {
		return errors.Trace(err)
	}

	if err := v.guest.sysVol.CopyFrom(v.guest.Image); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (v *virtGuest) allocGuest() error {
	if err := v.dom.define(); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func (v *virtGuest) domain() *domain {
	return &domain{
		guest: v.guest,
		virt:  v.virt,
	}
}

func (v *virtGuest) newFlock() *util.Flock {
	var fn = fmt.Sprintf("%s.flock", v.guest.Name())
	var fpth = filepath.Join(config.Conf.VirtFlockDir, fn)
	return util.NewFlock(fpth)
}

func (v *virtGuest) setNics() error {
	var leng = time.Duration(len(v.guest.nics))
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute*leng)
	defer cancel()

	for i, nic := range v.guest.nics {
		var dev = fmt.Sprintf("eth%d", i)

		var st = <-v.ga.Exec(ctx, "ip", "a", "add", nic.IP(), "dev", dev)
		if st.Error() != nil {
			return errors.Annotatef(st.Error(), "failed to `ip a add %s dev %s`", nic.IP(), dev)
		}
	}

	return nil
}
