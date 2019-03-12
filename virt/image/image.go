package image

import (
	"fmt"

	"github.com/projecteru2/yavirt/db"
	"github.com/projecteru2/yavirt/virt/common"
)

const (
	ImgUser = "user"

	selectQuery = `
SELECT id, parent_id, size, image_name, host_id, state, transit_status, create_time,
       transit_time, update_time
FROM image_tab`
)

type Image struct {
	common.Resource

	ParentID int64 `db:"parent_id"`
	Size     int
	Name     string `db:"image_name"`
	HostID   int64  `db:"host_id"`
}

func Load(id int64) (*Image, error) {
	var img = &Image{}
	var query = fmt.Sprintf("%s WHERE id=?", selectQuery)
	var err = db.Get(img, query, id)
	return img, err
}

func LoadByName(name string) (*Image, error) {
	var img = &Image{}
	var query = fmt.Sprintf("%s WHERE image_name=?", selectQuery)
	var err = db.Get(img, query, name)
	return img, err
}
