package netx

import (
	"testing"

	"github.com/projecteru2/yavirt/test/assert"
)

func TestIntToIPv4(t *testing.T) {
	assert.Equal(t, "192.168.1.1", IntToIPv4(3232235777))
	assert.Equal(t, "10.1.2.3", IntToIPv4(167838211))
	assert.Equal(t, "127.0.0.1", IntToIPv4(2130706433))
	assert.Equal(t, "255.255.255.255", IntToIPv4(4294967295))
}

func TestIPv4ToInt(t *testing.T) {
	var cases = []struct {
		out int64
		in  string
	}{
		{3232235777, "192.168.1.1"},
		{167838211, "10.1.2.3"},
		{2130706433, "127.0.0.1"},
		{4294967295, "255.255.255.255"},
	}

	for _, c := range cases {
		var i, err = IPv4ToInt(c.in)
		assert.NilErr(t, err)
		assert.Equal(t, c.out, i)
	}
}
