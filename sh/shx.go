package sh

import (
	"context"
	"io/ioutil"
	"os/exec"

	"github.com/juju/errors"
)

type shx struct{}

func (s shx) Remove(fpth string) error {
	return s.Exec(context.Background(), "rm", "-rf", fpth)
}

func (s shx) Copy(src, dest string) error {
	return s.Exec(context.Background(), "cp", src, dest)
}

func (s shx) Exec(ctx context.Context, name string, args ...string) error {
	var cmd = exec.CommandContext(ctx, name, args...)

	stderr, err := cmd.StderrPipe()
	if err != nil {
		return errors.Trace(err)
	}

	if err := cmd.Start(); err != nil {
		return errors.Trace(err)
	}

	slurp, err := ioutil.ReadAll(stderr)
	if err != nil {
		return errors.Trace(err)
	}

	if err := cmd.Wait(); err != nil {
		return errors.Annotatef(err, string(slurp))
	}

	return nil
}
