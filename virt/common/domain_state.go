package common

import "github.com/libvirt/libvirt-go"

type DomainState libvirt.DomainState

func (d DomainState) String() string {
	switch libvirt.DomainState(d) {
	case libvirt.DOMAIN_RUNNING:
		return "running"

	case libvirt.DOMAIN_BLOCKED:
		return "blocked"

	case libvirt.DOMAIN_PAUSED:
		return "paused"

	case libvirt.DOMAIN_SHUTDOWN:
		return "shutdowning"

	case libvirt.DOMAIN_CRASHED:
		return "crashed"

	case libvirt.DOMAIN_PMSUSPENDED:
		return "pmsuspended"

	case libvirt.DOMAIN_SHUTOFF:
		return "shutoff"

	case libvirt.DOMAIN_NOSTATE:
		fallthrough
	default:
		return "nostate"
	}
}
