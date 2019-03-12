package common

import (
	"fmt"
	"path/filepath"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/db"
	"github.com/projecteru2/yavirt/log"
)

const (
	TableGuest  = "guest_tab"
	TableVolume = "volume_tab"
)

type ResourceDomain interface {
	Filepath() string
}

type Resource struct {
	ID            int64
	Status        string `db:"state"`
	TransitStatus string `db:"transit_status"`
	CreateTime    int64  `db:"create_time"`
	TransitTime   int64  `db:"transit_time"`
	UpdateTime    int64  `db:"update_time"`
}

func (r Resource) JoinVirtPath(elem string) string {
	return filepath.Join(config.Conf.VirtDir, elem)
}

func (r *Resource) UpdateStatus(table, status string) error {
	if !r.CheckForwardStatus(status) {
		return NewForwardStatusErr(table, r.Status, status)
	}

	var query = fmt.Sprintf("UPDATE %s SET state=? WHERE id=?", table)
	var res, err = db.Exec(query, status, r.ID)
	if err != nil {
		return errors.Trace(err)
	}

	switch affect, err := res.RowsAffected(); {
	case err != nil:
		return errors.Trace(err)
	case affect != 1:
		log.Warnf("%s effected %d rows", query, affect)
	}

	r.Status = status

	return nil
}

func (r Resource) CheckForwardStatus(st string) bool {
	return checkForward(r.Status, st)
}
