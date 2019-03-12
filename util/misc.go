package util

import "github.com/juju/errors"

func Invoke(funcs []func() error) error {
	for _, fn := range funcs {
		if err := fn(); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}
