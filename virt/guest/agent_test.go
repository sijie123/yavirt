package guest

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"testing"

	"github.com/projecteru2/yavirt/test/assert"
)

func TestAgent(t *testing.T) {
	var agent = NewAgent("/tmp/virt/sock/guest-000001.sock")
	var in = "ping"
	var out = []byte("pong")

	var ret = ExecStatus{
		Exited:    true,
		Base64Out: base64.StdEncoding.EncodeToString(out),
	}

	enc, err := json.Marshal(ret)
	assert.NilErr(t, err)

	var qmp = &mockQmp{}
	qmp.On("Exec", in, []interface{}(nil)).Return([]byte(`{"pid":6735}`), nil)
	qmp.On("ExecStatus", 6735).Return(enc, nil)

	agent.qmp = qmp

	var st = <-agent.ExecOutput(context.Background(), in)
	assert.NotNil(t, st)
	assert.NilErr(t, st.Error())
	assert.Equal(t, 0, st.Code)

	so, err := st.Stdout()
	assert.NilErr(t, err)
	assert.Equal(t, out, so)

	se, err := st.Stderr()
	assert.NilErr(t, err)
	assert.Equal(t, []byte{}, se)

	qmp.AssertExpectations(t)
}

func TestRealAgent(t *testing.T) {
	// var agent = NewAgent("/tmp/virt/sock/guest-000001.sock")
	// assert.NilErr(t, agent.Ping(context.Background()))

	// var st = <-agent.ExecOutput(context.Background(), "ls", "-l", "/")
	// assert.NotNil(t, st)
	// assert.NilErr(t, st.Error())

	// so, err := st.Stdout()
	// assert.NilErr(t, err)
	// t.Logf("%s", so)

	// se, err := st.Stderr()
	// assert.NilErr(t, err)
	// assert.Equal(t, 0, len(se))
}
