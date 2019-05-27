package host

import (
	"fmt"

	"github.com/projecteru2/yavirt/db"
)

const selectQuery = "SELECT id, hostname, host_type, subnet, state, cpu, mem FROM host_tab"

type Host struct {
	ID     int64
	Host   string `db:"hostname"`
	Type   string `db:"host_type"`
	Subnet int64
	Status string `db:"state"`
	Cpu    int
	Mem    int64
}

func LoadByHost(hn string) (*Host, error) {
	var host = &Host{}
	var query = fmt.Sprintf("%s WHERE hostname=?", selectQuery)
	var err = db.Get(host, query, hn)
	return host, err
}

func Load(id int64) (*Host, error) {
	var host = &Host{}
	var query = fmt.Sprintf("%s WHERE id=?", selectQuery)
	var err = db.Get(host, query, id)
	return host, err
}
