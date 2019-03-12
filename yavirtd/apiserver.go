package yavirtd

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/api/types"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/ver"
)

func newApiHandler(yav *Server) http.Handler {
	gin.SetMode(gin.ReleaseMode)

	var srv = &apiServer{yav}
	var router = gin.Default()

	var v1 = router.Group("/v1")
	{
		v1.GET("/ping", srv.Ping)
		v1.GET("/info", srv.Info)
		v1.GET("/guests/:id", srv.GetGuest)
		v1.POST("/guests", srv.CreateGuest)
		v1.POST("/guests/stop", srv.StopGuest)
		v1.POST("/guests/start", srv.StartGuest)
		v1.POST("/guests/destroy", srv.DestroyGuest)
	}

	return router
}

type apiServer struct {
	yav *Server
}

func (s *apiServer) Info(c *gin.Context) {
	s.renderOK(c, types.HostInfo{
		ID:  fmt.Sprintf("%d", s.yav.host.ID),
		Cpu: s.yav.host.Cpu,
		Mem: s.yav.host.Mem,
	})
}

func (s *apiServer) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"version": ver.Version(),
	})
}

func (s *apiServer) dispatchMsg(c *gin.Context, req interface{}, fn func() error) {
	s.dispatch(c, req, func() (interface{}, error) {
		return nil, fn()
	})
}

type operate func() (interface{}, error)

func (s *apiServer) dispatch(c *gin.Context, req interface{}, fn operate) {
	if err := s.bind(c, req); err != nil {
		s.renderErr(c, err)
		return
	}

	var resp, err = fn()
	if err != nil {
		s.renderErr(c, err)
		return
	}

	if resp == nil {
		s.renderOKMsg(c)
	} else {
		s.renderOK(c, resp)
	}
}

func (s *apiServer) bind(c *gin.Context, req interface{}) error {
	switch c.Request.Method {
	case http.MethodGet:
		return c.ShouldBindUri(req)

	case http.MethodPost:
		return c.ShouldBind(req)

	default:
		return errors.Errorf("invalid HTTP method: %s", c.Request.Method)
	}
}

var okMsg = types.NewMsg("ok")

func (s *apiServer) renderOKMsg(c *gin.Context) {
	s.renderOK(c, okMsg)
}

func (s *apiServer) renderOK(c *gin.Context, resp interface{}) {
	c.JSON(http.StatusOK, resp)
}

func (s *apiServer) renderErr(c *gin.Context, err error) {
	log.Errorf(errors.ErrorStack(err))
	c.JSON(http.StatusInternalServerError, err.Error())
}
