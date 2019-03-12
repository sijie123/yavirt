package client

import (
	"context"
	"testing"

	"github.com/projecteru2/yavirt/api/types"
	"github.com/projecteru2/yavirt/test/assert"
	"github.com/projecteru2/yavirt/util"
)

func TestRealReq(t *testing.T) {
	var cli, err = New("127.0.0.1:9696", "v1")
	assert.NilErr(t, err)
	assert.NotNil(t, cli)

	// var g = testCreate(t, cli)
	// var id = g.ID
	// t.Logf("=== id: %s ===", id)

	// testStop(t, cli, id)
	// testStart(t, cli, id)
	// testStop(t, cli, id)
	// testDestroy(t, cli, id)
}

func testStart(t *testing.T, cli *Client, id string) {
	testReq(t, cli, id, cli.StartGuest)
}

func testDestroy(t *testing.T, cli *Client, id string) {
	testReq(t, cli, id, cli.DestroyGuest)
}

func testStop(t *testing.T, cli *Client, id string) {
	testReq(t, cli, id, cli.StopGuest)
}

func testReq(t *testing.T, cli *Client, id string, fn func(context.Context, string) (types.Msg, error)) {
	var resp, err = fn(nil, id)
	assert.NilErr(t, err)
	assert.NotNil(t, resp)
}

func testCreate(t *testing.T, cli *Client) types.Guest {
	var req = types.CreateGuestReq{Cpu: 1, Mem: util.GB, ImageID: 1}
	var resp, err = cli.CreateGuest(nil, req)
	assert.NilErr(t, err)
	assert.NotNil(t, resp)
	return resp
}
