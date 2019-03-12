package common

import (
	"fmt"

	"github.com/libvirt/libvirt-go"
)

type DomainStateErr struct {
	Resource string
	Exp      DomainState
	Act      DomainState
}

func NewDomainStateErr(res string, exp, act libvirt.DomainState) DomainStateErr {
	return DomainStateErr{
		Resource: res,
		Exp:      DomainState(exp),
		Act:      DomainState(act),
	}
}

func (e DomainStateErr) Error() string {
	return fmt.Sprintf("require %s, but %s is %s", e.Exp, e.Resource, e.Act)
}

type ForwardStatusErr struct {
	res  string
	src  string
	dest string
}

func NewForwardStatusErr(res, src, dest string) ForwardStatusErr {
	return ForwardStatusErr{
		res:  res,
		src:  src,
		dest: dest,
	}
}

func (e ForwardStatusErr) Error() string {
	return fmt.Sprintf("cannot forward %s %s => %s", e.res, e.src, e.dest)
}
