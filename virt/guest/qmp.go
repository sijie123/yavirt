package guest

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/util"
)

const QmpGuestExecStatus = "guest-exec-status"

type qmp struct {
	sync.Mutex

	// whether is guest-agent which means virsh qemu-agent-command,
	// the false value indicates virsh qemu-monitor-command.
	ga bool

	sockfile string
	sock     net.Conn
	reader   *bufio.Reader
	writer   *bufio.Writer

	greeting *json.RawMessage
}

type qmpResp struct {
	Event    *json.RawMessage
	Greeting *json.RawMessage
	Return   *json.RawMessage
	Error    *qmpError
}

type qmpError struct {
	Class string
	Desc  string
}

func (e *qmpError) Error() string {
	return fmt.Sprintf("QMP error %s: %s", e.Class, e.Desc)
}

func newQmp(sockfile string, ga bool) *qmp {
	return &qmp{
		sockfile: sockfile,
		ga:       ga,
	}
}

func (q *qmp) Exec(path string, args []interface{}, output bool) ([]byte, error) {
	q.Lock()
	defer q.Unlock()

	if err := q.connect(); err != nil {
		return nil, errors.Trace(err)
	}

	var exArg = map[string]interface{}{
		"path":           path,
		"capture-output": output,
	}
	if args != nil {
		exArg["arg"] = args
	}

	var buf, err = newQmpCmd("guest-exec", exArg).bytes()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return q.exec(buf)
}

func (q *qmp) ExecStatus(pid int) ([]byte, error) {
	q.Lock()
	defer q.Unlock()

	var buf, err = newQmpCmd("guest-exec-status", map[string]interface{}{"pid": pid}).bytes()
	if err != nil {
		return nil, errors.Trace(err)
	}

	return q.exec(buf)
}

func (q *qmp) exec(cmd []byte) ([]byte, error) {
	switch resp, err := q.req(cmd); {
	case err != nil:
		return nil, errors.Trace(err)

	case resp.Error != nil:
		return nil, errors.Trace(resp.Error)

	default:
		return []byte(*resp.Return), nil
	}
}

func (q *qmp) connect() error {
	if q.sock != nil {
		return nil
	}

	var sock, err = net.DialTimeout("unix", q.sockfile, time.Second*8)
	if err != nil {
		return errors.Trace(err)
	}

	q.sock = sock
	q.reader = bufio.NewReader(q.sock)
	q.writer = bufio.NewWriter(q.sock)

	if !q.ga {
		if err := q.handshake(); err != nil {
			q.Close()
			return errors.Trace(err)
		}
	}

	return nil
}

func (q *qmp) handshake() error {
	return util.Invoke([]func() error{
		q.greet,
		q.capabilities,
	})
}

func (q *qmp) capabilities() error {
	var cmd, err = newQmpCmd("qmp_capabilities", nil).bytes()
	if err != nil {
		return errors.Trace(err)
	}

	switch resp, err := q.req(cmd); {
	case err != nil:
		return errors.Trace(err)

	case resp.Return == nil:
		return errors.Errorf("QMP negotiation error")

	default:
		return nil
	}
}

func (q *qmp) greet() error {
	var buf, err = q.read()
	if err != nil {
		return errors.Trace(err)
	}

	var resp qmpResp

	switch err := json.Unmarshal(buf, &resp.Greeting); {
	case err != nil:
		return errors.Trace(err)
	case resp.Greeting == nil:
		return errors.Errorf("QMP greeting error")
	}

	q.greeting = resp.Greeting

	return nil
}

func (q *qmp) Close() (err error) {
	if q.sock != nil {
		err = q.sock.Close()
	}
	return
}

func (q *qmp) req(cmd []byte) (qmpResp, error) {
	var resp qmpResp

	if err := q.write(cmd); err != nil {
		return resp, errors.Trace(err)
	}

	var buf, err = q.read()
	if err != nil {
		return resp, errors.Trace(err)
	}

	if err := json.Unmarshal(buf, &resp); err != nil {
		return resp, errors.Trace(err)
	}

	return resp, nil
}

func (q *qmp) write(buf []byte) error {
	if _, err := q.writer.Write(append(buf, '\x0a')); err != nil {
		return errors.Trace(err)
	}

	if err := q.writer.Flush(); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (q *qmp) read() ([]byte, error) {
	for {
		var buf, err = q.reader.ReadBytes('\n')
		if err != nil {
			return nil, errors.Trace(err)
		}

		var resp qmpResp
		if err := json.Unmarshal(buf, &resp); err != nil {
			return nil, errors.Trace(err)
		}

		if resp.Event != nil {
			log.Infof("recv event: %s", resp.Event)
			continue
		}

		return buf, nil
	}
}

type qmpCmd struct {
	Name string                 `json:"execute"`
	Args map[string]interface{} `json:"arguments,omitempty"`
}

func newQmpCmd(name string, args map[string]interface{}) (c qmpCmd) {
	c.Name = name
	c.Args = args
	return
}

func (c qmpCmd) bytes() ([]byte, error) {
	return json.Marshal(c)
}
