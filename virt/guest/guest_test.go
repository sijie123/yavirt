package guest

import (
	_ "os"
	"testing"

	"github.com/projecteru2/yavirt/test/assert"
	"github.com/projecteru2/yavirt/test/mock"
	_ "github.com/projecteru2/yavirt/util"
	"github.com/projecteru2/yavirt/virt/common"
	"github.com/projecteru2/yavirt/virt/host"
	"github.com/projecteru2/yavirt/virt/image"
	"github.com/projecteru2/yavirt/virt/nic"
	"github.com/projecteru2/yavirt/virt/test"
	"github.com/projecteru2/yavirt/virt/volume"
)

type mockVolumesOp struct {
	mock.Mock
	vols   []*volume.Volume
	sysVol *volume.Volume
}

func newMockVolumesOp() *mockVolumesOp {
	var vol = &volume.Volume{}
	vol.ID = 1

	var volsOp = &mockVolumesOp{
		sysVol: vol,
		vols:   []*volume.Volume{vol},
	}

	return volsOp
}

func (m *mockVolumesOp) SysVolume() (*volume.Volume, error) {
	return m.sysVol, nil
}

func (m *mockVolumesOp) Append(vol *volume.Volume) {
	m.vols = append(m.vols, vol)
}

func (m *mockVolumesOp) UpdateStatus(st string) error {
	for _, v := range m.vols {
		v.Status = st
	}
	return nil
}

func TestLifecycle(t *testing.T) {
	var cancel = test.MockAll()
	defer cancel()

	var guest, bot = newMockedGuest()
	guest.host = &host.Host{ID: 1, Subnet: 12625921}

	assert.NilErr(t, guest.Insert())

	assert.Equal(t, common.StatusPending, guest.Status)
	assert.Equal(t, common.StatusPending, guest.sysVol.Status)

	bot.On("Create").Return(nil)
	assert.NilErr(t, guest.Create())
	assert.Equal(t, common.StatusCreating, guest.Status)
	assert.Equal(t, common.StatusCreating, guest.sysVol.Status)

	bot.On("Boot").Return(nil)
	assert.NilErr(t, guest.Start())
	assert.Equal(t, common.StatusRunning, guest.Status)
	assert.Equal(t, common.StatusRunning, guest.sysVol.Status)

	bot.On("Shutdown").Return(nil)
	assert.NilErr(t, guest.Stop())
	assert.Equal(t, common.StatusStopped, guest.Status)
	assert.Equal(t, common.StatusStopped, guest.sysVol.Status)

	assert.NilErr(t, guest.Resize())
	assert.Equal(t, common.StatusStopped, guest.Status)
	assert.Equal(t, common.StatusStopped, guest.sysVol.Status)

	assert.NilErr(t, guest.Start())
	assert.Equal(t, common.StatusRunning, guest.Status)
	assert.Equal(t, common.StatusRunning, guest.sysVol.Status)

	assert.NilErr(t, guest.Stop())
	assert.Equal(t, common.StatusStopped, guest.Status)
	assert.Equal(t, common.StatusStopped, guest.sysVol.Status)

	bot.On("Undefine").Return(nil)
	assert.NilErr(t, guest.Destroy())
	assert.Equal(t, common.StatusDestroyed, guest.Status)
	assert.Equal(t, common.StatusDestroyed, guest.sysVol.Status)

	bot.AssertExpectations(t)
}

func TestLifecycle_InvalidStatus(t *testing.T) {
	var cancel = test.MockAll()
	defer cancel()

	var guest = &Guest{
		newBot: newMockBot,
		vols:   newMockVolumesOp(),
		Image:  &image.Image{},
		nics:   []*nic.Nic{},
	}

	guest.Status = common.StatusDestroyed
	assert.Err(t, guest.Insert())
	assert.Err(t, guest.Create())
	assert.Err(t, guest.Stop())
	assert.Err(t, guest.Start())

	guest.Status = common.StatusResizing
	assert.Err(t, guest.Destroy())

	guest.Status = common.StatusPending
	assert.Err(t, guest.Resize())
}

func TestSyncState(t *testing.T) {
	var cancel = test.MockAll()
	defer cancel()

	var guest, bot = newMockedGuest()

	guest.Status = common.StatusCreating
	bot.On("Create").Return(nil)
	bot.On("Boot").Return(nil)
	assert.NilErr(t, guest.SyncState())

	guest.Status = common.StatusDestroying
	bot.On("Undefine").Return(nil)
	assert.NilErr(t, guest.SyncState())

	guest.Status = common.StatusStopping
	bot.On("Shutdown").Return(nil)
	assert.NilErr(t, guest.SyncState())

	guest.Status = common.StatusStarting
	assert.NilErr(t, guest.SyncState())

	bot.AssertExpectations(t)
}

func newMockedGuest() (*Guest, *mockBot) {
	var bot, _ = newMockBot(nil)
	var mbot = bot.(*mockBot)

	var vols = newMockVolumesOp()

	var guest = &Guest{
		newBot: func(g *Guest) (Bot, error) {
			mbot.guest = g
			return mbot, nil
		},
		vols:  vols,
		Image: &image.Image{},
	}

	return guest, mbot
}

func TestRealGuest(t *testing.T) {
	// hn, err := os.Hostname()
	// assert.NilErr(t, err)

	// phy, err := host.LoadByHost(hn)
	// assert.NilErr(t, err)

	// guest, err := Create(1, util.GB, 1, phy)
	// assert.Nil(t, err)
	// assert.NotNil(t, guest)

	// gid := guest.ID
	// //gid := 1
	// assert.NilErr(t, Start(gid))
	// assert.NilErr(t, Start(gid))
}
