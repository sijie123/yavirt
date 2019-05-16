package nic

import (
	"context"
	"fmt"
	"time"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/db"
	"github.com/projecteru2/yavirt/netx"
	"github.com/projecteru2/yavirt/virt/common"
)

const selectQuery = "SELECT id, high_value, low_value, prefix, gateway, guest_id, addr_type, state, host_id FROM addr_tab"

func LoadGuestNics(guestID int64) ([]*Nic, error) {
	var query = fmt.Sprintf("%s WHERE guest_id=?", selectQuery)

	var nics = []*Nic{}

	if err := db.Select(&nics, query, guestID); err != nil {
		return nil, errors.Trace(err)
	}

	return nics, nil
}

func Alloc(guestID, hostID int64) (*Nic, error) {
	var ctx, cancel = context.WithTimeout(context.Background(), time.Minute)
	defer cancel()

	var nic = &Nic{}

	var selQuery = fmt.Sprintf("%s WHERE state='%s' LIMIT 1 FOR UPDATE",
		selectQuery, common.AddrStatusFree)

	var err = db.SelectForUpdate(ctx, nic, selQuery, func(obj interface{}) string {
		return fmt.Sprintf(`UPDATE addr_tab
SET guest_id=%d, state='%s', host_id=%d
WHERE id=%d AND state='%s'`,
			guestID,
			common.AddrStatusOccupied,
			hostID,
			nic.ID,
			common.AddrStatusFree,
		)
	})

	if err != nil {
		return nil, errors.Trace(err)
	}

	return nic, nil
}

type Nic struct {
	ID        int64
	HighValue int64 `db:"high_value"`
	LowValue  int64 `db:"low_value"`
	Prefix    int
	Gateway   string
	GuestID   int64  `db:"guest_id"`
	Type      string `db:"addr_type"`
	Status    string `db:"state"`
	HostID    int64  `db:"host_id"`
}

func (n *Nic) Free() (err error) {
	_, err = db.Exec("UPDATE addr_tab SET guest_id=0, state=?, host_id=0 WHERE id=?",
		common.AddrStatusFree, n.ID)
	return
}

func (n *Nic) IP() string {
	return fmt.Sprintf("%s/%d", netx.IntToIPv4(n.LowValue), n.Prefix)
}
