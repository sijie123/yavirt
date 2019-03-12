package guest

import (
	"path/filepath"
	"time"

	"github.com/juju/errors"
	"github.com/libvirt/libvirt-go"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/util"
	"github.com/projecteru2/yavirt/virt/common"
	"github.com/projecteru2/yavirt/virt/template"
)

const ConnectListAllDomainsFlags = libvirt.CONNECT_LIST_DOMAINS_ACTIVE |
	libvirt.CONNECT_LIST_DOMAINS_INACTIVE |
	libvirt.CONNECT_LIST_DOMAINS_PERSISTENT |
	libvirt.CONNECT_LIST_DOMAINS_TRANSIENT |
	libvirt.CONNECT_LIST_DOMAINS_RUNNING |
	libvirt.CONNECT_LIST_DOMAINS_PAUSED |
	libvirt.CONNECT_LIST_DOMAINS_SHUTOFF |
	libvirt.CONNECT_LIST_DOMAINS_OTHER |
	libvirt.CONNECT_LIST_DOMAINS_MANAGEDSAVE |
	libvirt.CONNECT_LIST_DOMAINS_NO_MANAGEDSAVE |
	libvirt.CONNECT_LIST_DOMAINS_AUTOSTART |
	libvirt.CONNECT_LIST_DOMAINS_NO_AUTOSTART |
	libvirt.CONNECT_LIST_DOMAINS_HAS_SNAPSHOT |
	libvirt.CONNECT_LIST_DOMAINS_NO_SNAPSHOT

type domain struct {
	dom   *libvirt.Domain
	guest *Guest
	virt  *libvirt.Connect
}

func (d *domain) define() error {
	var buf, err = d.render()
	if err != nil {
		return errors.Trace(err)
	}

	var create = func() (err error) {
		d.dom, err = d.virt.DomainDefineXMLFlags(string(buf), libvirt.DOMAIN_DEFINE_VALIDATE)
		return
	}

	for i := 0; ; i++ {
		time.Sleep(time.Second * time.Duration(i))
		i %= 5

		switch st, err := d.getState(); {
		case err != nil:
			if !d.isNoDomainErr(err) {
				return errors.Trace(err)
			}
			if err := create(); err != nil {
				return errors.Trace(err)
			}
			continue

		case st == libvirt.DOMAIN_RUNNING:
			// Force shuting down.
			if err := d.dom.Destroy(); err != nil {
				return errors.Trace(err)
			}
			// Reboot it.
			fallthrough

		case st == libvirt.DOMAIN_SHUTOFF:
			return d.boot()

		default:
			return common.NewDomainStateErr("guest", libvirt.DOMAIN_NOSTATE, st)
		}
	}
}

func (d *domain) boot() error {
	var expState = libvirt.DOMAIN_SHUTOFF

	for i := 0; ; i++ {
		time.Sleep(time.Second * time.Duration(i))
		i %= 5

		switch st, err := d.getState(); {
		case err != nil:
			return errors.Trace(err)

		case st == libvirt.DOMAIN_RUNNING:
			return nil

		case st == expState:
			if err := d.dom.Create(); err != nil {
				return errors.Trace(err)
			}
			continue

		default:
			return common.NewDomainStateErr("guest", expState, st)
		}
	}
}

func (d *domain) shutdown() error {
	var expState = libvirt.DOMAIN_RUNNING

	for i := 0; ; i++ {
		time.Sleep(time.Second * time.Duration(i))
		i %= 5

		switch st, err := d.getState(); {
		case err != nil:
			return errors.Trace(err)

		case st == libvirt.DOMAIN_SHUTOFF:
			return nil

		case st == libvirt.DOMAIN_SHUTDOWN:
			// It's shutting now, wait to be shutoff.
			continue

		case st == expState:
			if err := d.dom.ShutdownFlags(libvirt.DOMAIN_SHUTDOWN_DEFAULT); err != nil {
				return errors.Trace(err)
			}
			continue

		default:
			return common.NewDomainStateErr("guest", expState, st)
		}
	}
}

func (d *domain) undefine() error {
	var expState = libvirt.DOMAIN_SHUTOFF

	switch st, err := d.getState(); {
	case err != nil:
		if d.isNoDomainErr(err) {
			return nil
		}
		return errors.Trace(err)

	case st == expState:
		return d.dom.UndefineFlags(libvirt.DOMAIN_UNDEFINE_MANAGED_SAVE)

	default:
		return common.NewDomainStateErr("guest", expState, st)
	}
}

func (d *domain) isNoDomainErr(err error) bool {
	switch raw, ok := errors.Cause(err).(libvirt.Error); {
	case ok:
		return raw.Code == libvirt.ERR_NO_DOMAIN
	default:
		return false
	}
}

func (d *domain) close() (err error) {
	if d.dom != nil {
		err = d.dom.Free()
	}
	return
}

func (d *domain) getState() (st libvirt.DomainState, err error) {
	if err := d.load(); err != nil {
		return libvirt.DOMAIN_NOSTATE, errors.Trace(err)
	}

	st, _, err = d.dom.GetState()

	return
}

func (d *domain) load() (err error) {
	d.dom, err = d.virt.LookupDomainByName(d.guest.Name())
	return
}

func (d *domain) render() ([]byte, error) {
	uuid, err := util.UuidStr()
	if err != nil {
		return nil, errors.Trace(err)
	}

	var args = map[string]interface{}{
		"name":       d.name(),
		"uuid":       uuid,
		"memory":     d.guest.convMem2MB(),
		"cpu":        d.guest.Cpu,
		"sysvol":     d.guest.sysVol.Filepath(),
		"gasock":     d.guest.sockfile(),
		"virtbridge": config.Conf.VirtBridge,
	}

	return template.Render(d.templateFilepath(), args)
}

func (d *domain) templateFilepath() string {
	return filepath.Join(config.Conf.VirtTmplDir, "guest.xml")
}

func (d *domain) name() string {
	return d.guest.Name()
}
