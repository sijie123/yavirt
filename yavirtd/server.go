package yavirtd

import (
	"context"
	"net"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/metric"
	"github.com/projecteru2/yavirt/netx"
	"github.com/projecteru2/yavirt/virt/host"
)

type Server struct {
	lapi       net.Listener
	addr       string
	httpServer *http.Server
	exit       struct {
		sync.Once
		ch chan struct{}
	}

	host *host.Host
}

func New() (*Server, error) {
	var srv = &Server{}
	srv.exit.ch = make(chan struct{}, 1)

	if err := srv.setup(); err != nil {
		return nil, errors.Trace(err)
	}

	return srv, nil
}

func (s *Server) setup() (err error) {
	if s.lapi, s.addr, err = s.listen(config.Conf.BindAddr); err != nil {
		return errors.Trace(err)
	}

	s.httpServer = s.newHttpServer()

	if err = s.setupHost(); err != nil {
		return errors.Trace(err)
	}

	return
}

func (s *Server) setupHost() error {
	hn, err := os.Hostname()
	if err != nil {
		return errors.Trace(err)
	}

	cur, err := host.LoadByHost(hn)
	if err != nil {
		return errors.Trace(err)
	}

	s.host = cur

	return nil
}

func (s *Server) newHttpServer() *http.Server {
	var mux = http.NewServeMux()
	mux.Handle("/metrics", metric.Handler())
	mux.Handle("/", newApiHandler(s))
	return &http.Server{Handler: mux}
}

func (s *Server) listen(addr string) (lis net.Listener, ip string, err error) {
	var network = "tcp"
	if lis, err = net.Listen(network, addr); err != nil {
		return
	}

	if ip, err = netx.GetLocalIP(network, lis.Addr().String()); err != nil {
		return
	}

	return
}

func (s *Server) Reload() error {
	return nil
}

func (s *Server) Run() (err error) {
	defer func() {
		log.Warnf("[yavirtd] yavirtd server %p loop exit", s)
		s.Close()
	}()

	if err := s.disasterRecover(); err != nil {
		return errors.Trace(err)
	}

	var errCh = make(chan error, 1)
	go func() {
		defer func() {
			log.Warnf("[yavirtd] HTTP server %p exit", s.httpServer)
		}()
		errCh <- s.httpServer.Serve(s.lapi)
	}()

	select {
	case <-s.exit.ch:
		return nil
	case err = <-errCh:
		return errors.Trace(err)
	}
}

func (s *Server) Close() {
	s.exit.Do(func() {
		close(s.exit.ch)

		var ctx, cancel = context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()

		if err := s.httpServer.Shutdown(ctx); err != nil {
			log.Errorf(errors.ErrorStack(err))
			metric.IncrError()
		}
	})
}

func (s *Server) ExitCh() chan struct{} {
	return s.exit.ch
}
