package guest

import "github.com/projecteru2/yavirt/test/mock"

type mockQmp struct {
	mock.Mock
}

func (q *mockQmp) Exec(cmd string, args []interface{}, stdio bool) ([]byte, error) {
	var ret = mock.NewRet(q.Called(cmd, args))
	return ret.Bytes(0), ret.Err(1)
}

func (q *mockQmp) ExecStatus(pid int) ([]byte, error) {
	var ret = mock.NewRet(q.Called(pid))
	return ret.Bytes(0), ret.Err(1)
}

func (q *mockQmp) Close() error {
	var ret = mock.NewRet(q.Called())
	return ret.Err(0)
}
