package vnet

type Addr struct {
	ID     int64
	Type   string
	Status string
	HostID int64 `db:"host_id"`
}
