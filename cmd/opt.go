package cmd

import (
	"fmt"

	"github.com/docopt/docopt-go"
	"github.com/juju/errors"
)

var (
	ConfigFilesKey = "CONFIG-FILES"
	ConfigFilesOpt = fmt.Sprintf("[%s]...", ConfigFilesKey)
)

func Parse(usage, ver string) *Options {
	opt, err := parse(usage, ver)
	if err != nil {
		PanicError(err)
	}
	return opt
}

func parse(doc, ver string) (*Options, error) {
	var opt, err = docopt.Parse(doc, nil, true, ver, false)
	if err != nil {
		return nil, errors.Trace(err)
	}
	return &Options{opt}, nil
}

type Options struct {
	opts map[string]interface{}
}

func (opt *Options) MustBoolArg(name string) bool {
	var arg, err = opt.BoolArg(name)
	if err != nil {
		panic(errors.ErrorStack(err))
	}
	return arg
}

func (opt *Options) BoolArg(name string) (bool, error) {
	var arg, err = opt.getarg(name)
	if err != nil {
		return false, errors.Trace(err)
	}

	var ret, ok = arg.(bool)
	if !ok {
		return false, errors.Errorf("invalid bool arg: %s", name)
	}

	return ret, nil
}

func (opt *Options) MustStrings(name string) []string {
	var strs, err = opt.Strings(name)
	if err != nil {
		PanicError(err)
	}
	return strs
}

func (opt *Options) Strings(name string) ([]string, error) {
	var arg, err = opt.getarg(name)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var strs, ok = arg.([]string)
	if !ok {
		return nil, errors.Errorf("invalid str slice arg: %s", name)
	}

	return strs, nil
}

func (opt *Options) getarg(name string) (interface{}, error) {
	if opt.opts == nil {
		return nil, errors.Errorf("invalid options")
	}

	var arg, exists = opt.opts[name]
	if !exists {
		return nil, errors.Errorf("no such arg: %s", name)
	}

	return arg, nil
}
