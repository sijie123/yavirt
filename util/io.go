package util

import (
	"io/ioutil"
	"os"

	"github.com/juju/errors"
)

func ReadAll(fpth string) ([]byte, error) {
	f, err := os.Open(fpth)
	if err != nil {
		return nil, errors.Trace(err)
	}

	buf, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, errors.Trace(err)
	}

	return buf, nil
}

func WriteTempFile(buf []byte) (string, error) {
	f, err := ioutil.TempFile(os.TempDir(), "temp-guest-*.xml")
	if err != nil {
		return "", errors.Trace(err)
	}

	if _, err := f.Write(buf); err != nil {
		return "", errors.Trace(err)
	}

	return f.Name(), nil
}
