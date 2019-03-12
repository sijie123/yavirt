package cmd

import (
	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/log"
)

func Panic(fmt string, args ...interface{}) {
	PanicError(errors.Errorf(fmt, args...))
}

func PanicError(err error) {
	log.Fatalf(errors.ErrorStack(err))
}
