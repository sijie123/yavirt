package yavirtd

import (
	"strings"

	"github.com/projecteru2/yavirt/api/types"
	"github.com/projecteru2/yavirt/virt/guest"
)

func convGuestResp(g *guest.Guest) (resp types.Guest) {
	resp.ID = types.EruID(g.ID)
	resp.Status = g.Status
	resp.TransitStatus = g.TransitStatus
	resp.CreateTime = g.CreateTime
	resp.TransitTime = g.TransitTime
	resp.UpdateTime = g.UpdateTime
	resp.ImageID = g.ImageID
	resp.ImageName = g.Image.Name
	resp.Cpu = g.Cpu
	resp.Mem = g.Mem

	if ips := g.IPAddrs(); len(ips) > 0 {
		resp.Networks = map[string]string{"IP": strings.Join(ips, ", ")}
	}

	return
}
