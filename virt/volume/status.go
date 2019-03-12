package volume

import "github.com/projecteru2/yavirt/virt/common"

func (vol *Volume) setDestroying() error {
	return vol.UpdateStatus(common.StatusDestroying)
}

func (vol *Volume) setDestroyed() error {
	return vol.UpdateStatus(common.StatusDestroyed)
}

func (vol *Volume) UpdateStatus(st string) error {
	return vol.Resource.UpdateStatus(common.TableVolume, st)
}
