package yavirtd

import (
	"github.com/gin-gonic/gin"

	"github.com/projecteru2/yavirt/api/types"
	"github.com/projecteru2/yavirt/virt/guest"
	"github.com/projecteru2/yavirt/virt/image"
)

func (s *apiServer) GetGuest(c *gin.Context) {
	var req types.GuestReq

	s.dispatch(c, &req, func() (interface{}, error) {
		id, err := req.VirtID()
		if err != nil {
			return nil, err
		}

		vm, err := guest.Load(id)
		if err != nil {
			return nil, err
		}

		return convGuestResp(vm), nil
	})
}

func (s *apiServer) DestroyGuest(c *gin.Context) {
	var req types.GuestReq
	s.dispatchMsg(c, &req, func() error {
		var id, err = req.VirtID()
		if err != nil {
			return err
		}
		return guest.Destroy(id)
	})
}

func (s *apiServer) StopGuest(c *gin.Context) {
	var req types.GuestReq
	s.dispatchMsg(c, &req, func() error {
		var id, err = req.VirtID()
		if err != nil {
			return err
		}
		return guest.Stop(id)
	})
}

func (s *apiServer) StartGuest(c *gin.Context) {
	var req types.GuestReq
	s.dispatchMsg(c, &req, func() error {
		var id, err = req.VirtID()
		if err != nil {
			return err
		}
		return guest.Start(id)
	})
}

func (s *apiServer) CreateGuest(c *gin.Context) {
	var req types.CreateGuestReq

	s.dispatch(c, &req, func() (interface{}, error) {
		imgID, err := s.getImageID(req)
		if err != nil {
			return nil, err
		}

		vm, err := guest.Create(req.Cpu, req.Mem, imgID, s.yav.host.ID)
		if err != nil {
			return nil, err
		}

		return convGuestResp(vm), nil
	})
}

func (s *apiServer) getImageID(req types.CreateGuestReq) (int64, error) {
	if req.ImageID > 0 {
		return req.ImageID, nil
	}

	var img, err = image.LoadByName(req.ImageName)
	if err != nil {
		return 0, err
	}

	return img.ID, nil
}
