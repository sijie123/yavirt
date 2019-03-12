package guest

import (
	"fmt"
	"path/filepath"
	"strconv"
	"time"

	"github.com/juju/errors"
	"github.com/libvirt/libvirt-go"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/db"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/metric"
	"github.com/projecteru2/yavirt/util"
	"github.com/projecteru2/yavirt/virt/common"
	"github.com/projecteru2/yavirt/virt/image"
	"github.com/projecteru2/yavirt/virt/nic"
	"github.com/projecteru2/yavirt/virt/volume"
)

func New(cpu int, mem, imageID, hostID int64) (guest *Guest) {
	guest = &Guest{
		Cpu:     cpu,
		Mem:     mem,
		ImageID: imageID,
		HostID:  hostID,
		vols:    volume.NewVolumes(),
		nics:    []*nic.Nic{},
	}

	return
}

func Destroy(id int64) error {
	return ctrl(id, func(g *Guest) error {
		return g.Destroy()
	})
}

func Stop(id int64) error {
	return ctrl(id, func(g *Guest) error {
		return g.Stop()
	})
}

func Start(id int64) error {
	return ctrl(id, func(g *Guest) error {
		return g.Start()
	})
}

func ctrl(id int64, fn func(*Guest) error) error {
	var g, err = Load(id)
	if err != nil {
		return errors.Trace(err)
	}
	return fn(g)
}

func Create(cpu int, mem, imageID, hostID int64) (*Guest, error) {
	var guest = New(cpu, mem, imageID, hostID)

	var err = util.Invoke([]func() error{
		guest.LoadImage,
		guest.Insert,
		guest.Create,
	})

	if err != nil {
		log.Errorf(errors.ErrorStack(err))
		metric.IncrError()
	}

	return guest, err
}

func ListLocals() ([]*Guest, error) {
	var fake = &Guest{}
	bot, err := fake.bot()
	if err != nil {
		return nil, errors.Trace(err)
	}

	virt, ok := bot.(*virtGuest)
	if !ok {
		return nil, errors.Errorf("invalid *virtGuest")
	}

	defer virt.Close()

	doms, err := virt.virt.ListAllDomains(ConnectListAllDomainsFlags)
	if err != nil {
		return nil, errors.Trace(err)
	}

	var guests = make([]*Guest, len(doms))

	for i, d := range doms {
		if guests[i], err = buildGuest(d); err != nil {
			return nil, errors.Trace(err)
		}
	}

	return guests, nil
}

func buildGuest(dom libvirt.Domain) (guest *Guest, err error) {
	var id int64
	if id, err = parseGuestID(dom); err != nil {
		return nil, errors.Trace(err)
	}

	guest, err = Load(id)

	return
}

func parseGuestID(dom libvirt.Domain) (id int64, err error) {
	var name string
	if name, err = dom.GetName(); err != nil {
		return 0, errors.Trace(err)
	}

	var _, raw = util.PartRight(name, "-")

	id, err = strconv.ParseInt(raw, 10, 64)

	return
}

func Load(id int64) (*Guest, error) {
	var guest = &Guest{
		vols: volume.NewVolumes(),
	}

	var query = `SELECT id, image_id, host_id, cpu, mem, state, transit_status, create_time, transit_time, update_time
FROM guest_tab
WHERE id=?`

	if err := db.Get(guest, query, id); err != nil {
		return nil, errors.Annotatef(err, "query guest %d error", id)
	}

	if err := util.Invoke([]func() error{
		guest.LoadImage,
		guest.LoadVolumes,
		guest.LoadNics,
	}); err != nil {
		return nil, errors.Trace(err)
	}

	return guest, nil
}

type Guest struct {
	common.Resource

	ImageID int64 `db:"image_id"`
	HostID  int64 `db:"host_id"`
	Cpu     int
	Mem     int64

	Image *image.Image

	vols   volume.VolumesOp
	sysVol *volume.Volume

	nics []*nic.Nic

	newBot func(*Guest) (Bot, error)
}

func (g *Guest) SyncState() error {
	switch g.Status {
	case common.StatusDestroying:
		return g.destroy()

	case common.StatusStopping:
		return g.stop()

	case common.StatusStarting:
		return g.start()

	case common.StatusCreating:
		return g.create()

	default:
		return nil
	}
}

func (g *Guest) Start() error {
	return util.Invoke([]func() error{
		g.setStarting,
		g.start,
	})
}

func (g *Guest) start() error {
	return g.botOperate(func(bot Bot) error {
		return util.Invoke([]func() error{
			bot.Boot,
			g.setRunning,
		})
	})
}

func (g *Guest) Resize() error {
	return util.Invoke([]func() error{
		g.setResizing,

		g.Migrate,

		g.setStopped,
	})
}

func (g *Guest) Migrate() error {
	return util.Invoke([]func() error{
		g.setMigrating,
		g.migrate,
	})
}

func (g *Guest) migrate() error {
	return g.botOperate(func(bot Bot) error {
		return bot.Migrate()
	})
}

func (g *Guest) Insert() error {
	if !g.CheckForwardStatus(common.StatusPending) {
		return common.NewForwardStatusErr("guest", g.Status, common.StatusPending)
	}

	g.Status = common.StatusPending

	var inserts = []func() error{
		g.insert,
		g.insertVolume,
		g.insertNic,
	}

	if err := util.Invoke(inserts); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func (g *Guest) Create() error {
	return util.Invoke([]func() error{
		g.setCreating,
		g.create,
	})
}

func (g *Guest) create() error {
	return g.botOperate(func(bot Bot) error {
		return util.Invoke([]func() error{
			bot.Create,

			g.setStarting,
			bot.Boot,

			g.setRunning,
		})
	})
}

func (g *Guest) Stop() error {
	return util.Invoke([]func() error{
		g.setStopping,
		g.stop,
	})
}

func (g *Guest) stop() error {
	return g.botOperate(func(bot Bot) error {
		return util.Invoke([]func() error{
			bot.Shutdown,
			g.setStopped,
		})
	})
}

func (g *Guest) Destroy() error {
	return util.Invoke([]func() error{
		g.setDestroying,
		g.destroy,
	})
}

func (g *Guest) destroy() error {
	return g.botOperate(func(bot Bot) error {
		return util.Invoke([]func() error{
			bot.Undefine,
			g.deleteNics,
			g.setDestroyed,
		})
	})
}

func (g *Guest) botOperate(ops func(bot Bot) error) error {
	var bot, err = g.bot()
	if err != nil {
		return errors.Trace(err)
	}

	defer bot.Close()

	return ops(bot)
}

func (g *Guest) insertVolume() error {
	var vol = &volume.Volume{
		Format:   "qcow2",
		Capacity: g.Image.Size,
		Type:     volume.VolSys,
		HostID:   g.HostID,
	}
	vol.Status = common.StatusPending
	vol.CreateTime = time.Now().Unix()

	if err := vol.Insert(); err != nil {
		return errors.Trace(err)
	}

	var query = `INSERT INTO guest_volume_tab (guest_id, volume_id) VALUES (?, ?)`
	if _, err := db.Exec(query, g.ID, vol.ID); err != nil {
		return errors.Trace(err)
	}

	g.vols.Append(vol)

	g.sysVol = vol

	return nil
}

func (g *Guest) insertNic() error {
	var n, err = nic.Alloc(g.ID, g.HostID)
	if err != nil {
		return errors.Trace(err)
	}

	g.nics = append(g.nics, n)

	return nil
}

func (g *Guest) deleteNics() error {
	for _, n := range g.nics {
		if err := n.Free(); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (g *Guest) insert() error {
	var fields = []string{"image_id", "host_id", "cpu", "mem", "state", "create_time"}
	var res, err = db.Insert(g, "guest_tab", fields...)
	if err != nil {
		return errors.Trace(err)
	}

	var id int64
	if id, err = res.LastInsertId(); err != nil {
		return errors.Trace(err)
	}

	g.ID = id

	return nil
}

func (g *Guest) bot() (Bot, error) {
	if g.newBot != nil {
		return g.newBot(g)
	}
	return newVirtGuest(g)
}

func (g *Guest) Name() string {
	return fmt.Sprintf("guest-%06d", g.ID)
}

func (g *Guest) LoadImage() (err error) {
	g.Image, err = image.Load(g.ImageID)
	return
}

func (g *Guest) LoadNics() (err error) {
	g.nics, err = nic.LoadGuestNics(g.ID)
	return
}

func (g *Guest) LoadVolumes() (err error) {
	if g.vols, err = volume.LoadGuestVolumes(g.ID); err != nil {
		return errors.Trace(err)
	}

	g.sysVol, err = g.vols.SysVolume()

	return
}

func (g *Guest) sockfile() string {
	var fn = fmt.Sprintf("%s.sock", g.Name())
	return filepath.Join(config.Conf.VirtSockDir, fn)
}

func (g *Guest) convMem2MB() int64 {
	return util.ConvToMB(g.Mem)
}
