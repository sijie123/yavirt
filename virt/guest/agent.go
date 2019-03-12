package guest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"time"

	"github.com/juju/errors"
)

type Agent struct {
	qmp interface {
		Exec(cmd string, args []interface{}, stdio bool) ([]byte, error)
		ExecStatus(pid int) ([]byte, error)
		Close() error
	}
}

func NewAgent(sockfile string) *Agent {
	return &Agent{
		qmp: newQmp(sockfile, true),
	}
}

func (a *Agent) Ping(ctx context.Context) error {
	var st = <-a.Exec(ctx, "echo")
	return st.Error()
}

func (a *Agent) ExecOutput(ctx context.Context, cmd string, args ...interface{}) <-chan ExecStatus {
	return a.exec(ctx, cmd, args, true)
}

func (a *Agent) Exec(ctx context.Context, cmd string, args ...interface{}) <-chan ExecStatus {
	return a.exec(ctx, cmd, args, false)
}

func (a *Agent) exec(ctx context.Context, cmd string, args []interface{}, stdio bool) <-chan ExecStatus {
	var done = make(chan ExecStatus, 1)
	var st ExecStatus

	var data []byte
	data, st.err = a.qmp.Exec(cmd, args, stdio)
	if st.err != nil {
		done <- st
		return done
	}

	var ret = struct {
		Pid int
	}{}
	if st.err = a.decode(data, &ret); st.err != nil {
		done <- st
		return done
	}

	var pid = ret.Pid

	go func() {
		defer func() {
			done <- st
		}()

		var next = time.After(time.Millisecond)

		for i := 1; ; i++ {
			i %= 100

			select {
			case <-ctx.Done():
				st.err = errors.Annotatef(ctx.Err(), "exec %s error", cmd)
				return

			case <-next:
				if st = a.execStatus(pid, stdio); st.err != nil || st.Exited {
					return
				}
				next = time.After(time.Millisecond * time.Duration(i*10))
			}
		}
	}()

	return done
}

func (a *Agent) execStatus(pid int, stdio bool) (st ExecStatus) {
	var data, err = a.qmp.ExecStatus(pid)
	if err != nil {
		st.err = errors.Trace(err)
		return
	}

	if err := a.decode(data, &st); err != nil {
		st.err = errors.Trace(err)
	}

	return
}

func (a *Agent) Close() (err error) {
	if a.qmp != nil {
		err = a.qmp.Close()
	}
	return
}

func (a *Agent) decode(data []byte, v interface{}) error {
	return json.Unmarshal(data, v)
}

type ExecStatus struct {
	Exited       bool   `json:"exited"`
	Code         int    `json:"exitcode"`
	Base64Out    string `json:"out-data"`
	OutTruncated bool   `json:"out-truncated"`
	Base64Err    string `json:"err-data"`
	ErrTruncated bool   `json:"err-truncated"`

	err error
}

func (s ExecStatus) Stdout() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.Base64Out)
}

func (s ExecStatus) Stderr() ([]byte, error) {
	return base64.StdEncoding.DecodeString(s.Base64Err)
}

func (s ExecStatus) Error() error {
	switch {
	case s.err != nil:
		return errors.Trace(s.err)

	case !s.Exited:
		return errors.Errorf("still running")

	case s.Code != 0:
		return errors.Errorf("return %d; stdout: %s; stderr: %s", s.Code, s.Base64Out, s.Base64Err)

	default:
		return nil
	}
}
