package sh

import "context"

var shell Shell = shx{}

type Shell interface {
	Copy(src, dest string) error
	Remove(fpth string) error
	Exec(ctx context.Context, name string, args ...string) error
}

func Remove(fpth string) error {
	return shell.Remove(fpth)
}

func Copy(src, dest string) error {
	return shell.Copy(src, dest)
}

func ExecContext(ctx context.Context, name string, args ...string) error {
	return shell.Exec(ctx, name, args...)
}
