package main

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/cmd"
	"github.com/projecteru2/yavirt/config"
	"github.com/projecteru2/yavirt/log"
	"github.com/projecteru2/yavirt/metric"
	"github.com/projecteru2/yavirt/schema"
	"github.com/projecteru2/yavirt/ver"
	"github.com/projecteru2/yavirt/yavirtd"
)

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	var opt = cmd.Parse(usage(), ver.Version())
	if err := loadConfig(opt); err != nil {
		cmd.PanicError(err)
	}

	if err := setup(); err != nil {
		cmd.PanicError(err)
	}

	branch(opt)

	go cmd.Prof(config.Conf.ProfHttpPort)

	start()
}

func start() {
	var srv, err = yavirtd.New()
	if err != nil {
		cmd.PanicError(err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		handleSigns(srv)
	}()

	log.Infof("[main] server %p is running", srv)

	if err := srv.Run(); err != nil {
		log.Errorf(errors.ErrorStack(err))
		metric.IncrError()
	}

	wg.Wait()

	log.Warnf("[main] proc exit")
}

func branch(opt *cmd.Options) {
	switch {
	case opt.MustBoolArg("--init"):
		if err := schema.InitSchema(); err != nil {
			cmd.PanicError(err)
		}

	default:
		return
	}

	os.Exit(0)
}

func loadConfig(opt *cmd.Options) error {
	var files = opt.MustStrings(cmd.ConfigFilesKey)
	if err := config.Conf.Load(files); err != nil {
		return errors.Trace(err)
	}
	return nil
}

func setup() error {
	if err := log.Setup(config.Conf.LogLevel, config.Conf.LogFile); err != nil {
		return errors.Trace(err)
	}
	return nil
}

var signs = []os.Signal{
	syscall.SIGHUP,
	syscall.SIGINT,
	syscall.SIGTERM,
	syscall.SIGQUIT,
	syscall.SIGUSR2,
}

func handleSigns(srv *yavirtd.Server) {
	defer func() {
		log.Warnf("[main] signal handler exit")
		srv.Close()
	}()

	var signCh = make(chan os.Signal, 1)
	signal.Notify(signCh, signs...)

	var exit = srv.ExitCh()

	for {
		select {
		case sign := <-signCh:
			switch sign {
			case syscall.SIGUSR2:
				log.Warnf("[main] got sign USR2 to reload")
				srv.Reload()
			default:
				log.Warnf("[main] got sign %d to exit", sign)
				return
			}

		case <-exit:
			log.Warnf("[main] recv from server's exit ch")
			return
		}
	}
}

func usage() string {
	return fmt.Sprintf(`Usage:
    yavirtd %s
    yavirtd [--init | --version]`, cmd.ConfigFilesOpt)
}
