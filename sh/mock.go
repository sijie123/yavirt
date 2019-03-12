package sh

import (
	"context"

	"github.com/projecteru2/yavirt/test/mock"
)

type MockShell struct {
	mock.Mock
}

func NewMockShell() (Shell, func()) {
	var s = &MockShell{}
	return s, mockSh(s)
}

func mockSh(s Shell) func() {
	var old = shell

	shell = s

	return func() {
		shell = old
	}
}

func (s *MockShell) Remove(fpth string) error {
	return nil
}

func (s *MockShell) Copy(src, dest string) error {
	return nil
}

func (s *MockShell) Exec(ctx context.Context, name string, args ...string) error {
	return nil
}
