package guest

import (
	"testing"

	_ "github.com/libvirt/libvirt-go"

	"github.com/projecteru2/yavirt/netx"
	"github.com/projecteru2/yavirt/test/assert"
	"github.com/projecteru2/yavirt/test/mock"
	"github.com/projecteru2/yavirt/util"
	"github.com/projecteru2/yavirt/virt/common"
	"github.com/projecteru2/yavirt/virt/image"
	"github.com/projecteru2/yavirt/virt/nic"
	"github.com/projecteru2/yavirt/virt/volume"
)

type mockBot struct {
	mock.Mock
	guest *Guest
}

func newMockBot(guest *Guest) (Bot, error) {
	return &mockBot{guest: guest}, nil
}

func (b *mockBot) Close() error {
	return nil
}

func (b *mockBot) Create() error {
	var ret = mock.NewRet(b.Called())
	return ret.Err(0)
}

func (b *mockBot) Boot() error {
	var ret = mock.NewRet(b.Called())
	return ret.Err(0)
}

func (b *mockBot) Migrate() error {
	return nil
}

func (b *mockBot) Shutdown() error {
	var ret = mock.NewRet(b.Called())
	return ret.Err(0)
}

func (b *mockBot) Undefine() error {
	var ret = mock.NewRet(b.Called())
	return ret.Err(0)
}

func TestRealBot(t *testing.T) {
	assert.Nil(t, nil)

	var guest = New(1, util.MB*512, 1, 1)
	guest.ID = 1
	guest.Status = common.StatusPending
	guest.Image = &image.Image{Name: "centos7"}
	guest.Image.ID = 1

	ui32, err := netx.IPv4ToInt("192.168.1.1")
	assert.NilErr(t, err)

	guest.sysVol = &volume.Volume{}
	guest.sysVol.ID = 1
	guest.vols = volume.NewVolumes()
	guest.vols.Append(guest.sysVol)
	guest.nics = []*nic.Nic{
		&nic.Nic{
			LowValue: ui32,
			Prefix:   24,
		},
	}

	// bot, err := guest.bot()
	// assert.NilErr(t, err)
	// assert.NotNil(t, bot)

	// assert.NilErr(t, bot.Create())
	// assert.NilErr(t, bot.Create())
	// assert.NilErr(t, bot.Create())

	// assert.NilErr(t, bot.Shutdown())
	// assert.NilErr(t, bot.Shutdown())
	// assert.NilErr(t, bot.Shutdown())

	// assert.NilErr(t, bot.Boot())
	// assert.NilErr(t, bot.Boot())
	// assert.NilErr(t, bot.Boot())

	// assert.NilErr(t, bot.Shutdown())
	// assert.NilErr(t, bot.Shutdown())
	// assert.NilErr(t, bot.Shutdown())

	// assert.NilErr(t, bot.Undefine())
	// assert.NilErr(t, bot.Undefine())
	// assert.NilErr(t, bot.Undefine())
}
