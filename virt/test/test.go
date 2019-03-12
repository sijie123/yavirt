package test

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/projecteru2/yavirt/db"
	"github.com/projecteru2/yavirt/sh"
)

func init() {
	InitDir()
}

func MockAll() func() {
	var _, cancel1 = db.NewMockPool()
	var _, cancel2 = sh.NewMockShell()

	return func() {
		cancel1()
		cancel2()
	}
}

func InitDir() {
	var dir = filepath.Join(os.TempDir(), "virt/flock")
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		panic(err)
	}

	var tmplDir, err = filepath.Abs("../template")
	if err != nil {
		panic(err)
	}

	if err := os.Symlink(tmplDir, "/tmp/virt/template"); err != nil && !isErrExist(err) {
		panic(err)
	}
}

func isErrExist(err error) bool {
	return strings.HasSuffix(err.Error(), "file exists")
}
