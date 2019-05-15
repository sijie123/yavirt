package util

import (
	"io/ioutil"
	"path/filepath"

	"github.com/juju/errors"
)

func BaseFilename(fpth string) (fn string, ext string) {
	var base = filepath.Base(fpth)
	return PartRight(base, ".")
}

func AbsDir(fpth string) (string, error) {
	if filepath.IsAbs(fpth) {
		return fpth, nil
	}
	return filepath.Abs(fpth)
}

// Re-implements as filepath.Walk doesn't follow symlinks.
func Walk(root string, fn filepath.WalkFunc) error {
	var infos, err = ioutil.ReadDir(root)
	if err != nil {
		return errors.Trace(err)
	}

	for _, info := range infos {
		if err := fn(filepath.Join(root, info.Name()), info, nil); err != nil {
			return errors.Trace(err)
		}
	}

	return nil
}
