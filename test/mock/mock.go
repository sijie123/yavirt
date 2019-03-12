package mock

import testify "github.com/stretchr/testify/mock"

type Mock = testify.Mock

type Ret struct {
	testify.Arguments
}

func NewRet(args testify.Arguments) *Ret {
	return &Ret{args}
}

func (r *Ret) Err(index int) (err error) {
	if obj := r.Get(index); obj != nil {
		err = obj.(error)
	}
	return
}

func (r *Ret) Bytes(index int) (buf []byte) {
	if obj := r.Get(index); obj != nil {
		buf = obj.([]byte)
	}
	return
}
