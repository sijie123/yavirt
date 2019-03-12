package volume

import (
	"fmt"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/db"
	"github.com/projecteru2/yavirt/virt/common"
)

const (
	VolSys  = "sys"
	VolData = "data"
)

func LoadGuestVolumes(guestID int64) (VolumesOp, error) {
	var query = `
SELECT vol.id, format, capacity, volume_type, host_id, state, transit_status, create_time, transit_time, update_time
FROM volume_tab AS vol
INNER JOIN guest_volume_tab as gv
WHERE vol.id=gv.volume_id AND gv.guest_id=?`

	var vols = []*Volume{}

	if err := db.Select(&vols, query, guestID); err != nil {
		return nil, errors.Trace(err)
	}

	var v = NewVolumes()
	v.vols = vols

	return v, nil
}

type VolumesOp interface {
	SysVolume() (*Volume, error)
	Append(*Volume)
	UpdateStatus(string) error
}

type Volume struct {
	common.Resource

	Format   string
	Capacity int
	Type     string `db:"volume_type"`
	HostID   int64  `db:"host_id"`

	newBot func(*Volume) (Bot, error)
}

func (vol *Volume) Undefine() error {
	var bot, err = vol.bot()
	if err != nil {
		return errors.Trace(err)
	}

	defer bot.Close()

	return bot.Undefine()
}

func (vol *Volume) Insert() error {
	var fields = []string{"format", "capacity", "volume_type", "host_id", "state", "create_time"}
	var res, err = db.Insert(vol, "volume_tab", fields...)
	if err != nil {
		return errors.Trace(err)
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return errors.Trace(err)
	}

	vol.ID = id

	return nil
}

func (vol *Volume) Name() string {
	switch vol.Type {
	case VolSys, "":
		return fmt.Sprintf("sys-%06d.vol", vol.ID)
	case VolData:
		return fmt.Sprintf("dat-%06d.vol", vol.ID)
	default:
		return ""
	}
}

func (vol *Volume) CopyFrom(src common.ResourceDomain) error {
	var bot, err = vol.bot()
	if err != nil {
		return errors.Trace(err)
	}

	defer bot.Close()

	return bot.CopyFrom(src)
}

func (vol *Volume) bot() (Bot, error) {
	if vol.newBot != nil {
		return vol.newBot(vol)
	}
	return newVirtVol(vol)
}

func (vol *Volume) Filepath() string {
	return vol.JoinVirtPath(vol.Name())
}

type Volumes struct {
	vols []*Volume
}

func NewVolumes() *Volumes {
	return &Volumes{vols: []*Volume{}}
}

func (v *Volumes) SysVolume() (*Volume, error) {
	for _, vol := range v.vols {
		if vol.Type == VolSys {
			return vol, nil
		}
	}
	return nil, errors.Errorf("no sys volume")
}

func (v *Volumes) Append(vol *Volume) {
	v.vols = append(v.vols, vol)
}

func (v *Volumes) Vols() []*Volume {
	return v.vols
}

func (v *Volumes) UpdateStatus(st string) error {
	for _, v := range v.vols {
		if err := v.UpdateStatus(st); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}
