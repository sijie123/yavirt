package assert

import (
	"testing"

	"github.com/juju/errors"
	"github.com/stretchr/testify/require"
)

func NilErr(t *testing.T, err error) {
	Nil(t, err, errors.ErrorStack(err))
}

func Err(t *testing.T, err error) {
	NotNil(t, err, errors.ErrorStack(err))
}

func Nil(t *testing.T, obj interface{}, msgAndArgs ...interface{}) {
	require.Nil(t, obj, msgAndArgs...)
}

func NotNil(t *testing.T, obj interface{}, msgAndArgs ...interface{}) {
	require.NotNil(t, obj, msgAndArgs...)
}

func True(t *testing.T, b bool, msgAndArgs ...interface{}) {
	Equal(t, true, b, msgAndArgs...)
}

func False(t *testing.T, b bool, msgAndArgs ...interface{}) {
	Equal(t, false, b, msgAndArgs...)
}

func Equal(t *testing.T, exp, act interface{}, msgAndArgs ...interface{}) {
	require.Equal(t, exp, act, msgAndArgs...)
}
