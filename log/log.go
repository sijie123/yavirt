package log

import (
	"os"

	"github.com/juju/errors"
	"github.com/sirupsen/logrus"
)

func Setup(level, file string) error {
	if err := setupLevel(level); err != nil {
		return errors.Trace(err)
	}

	logrus.SetFormatter(&logrus.TextFormatter{FullTimestamp: true})

	if err := setupOutput(file); err != nil {
		return errors.Trace(err)
	}

	return nil
}

func setupOutput(file string) error {
	if len(file) < 1 {
		return nil
	}

	var f, err = os.OpenFile(file, os.O_APPEND|os.O_CREATE|os.O_RDWR, os.ModePerm)
	if err != nil {
		return errors.Trace(err)
	}

	logrus.SetOutput(f)

	return nil
}

func setupLevel(level string) error {
	if len(level) < 1 {
		return nil
	}

	var lv, err = logrus.ParseLevel(level)
	if err != nil {
		return errors.Trace(err)
	}

	logrus.SetLevel(lv)

	return nil
}

func Warnf(fmt string, args ...interface{}) {
	logrus.Warnf(fmt, args...)
}

func Errorf(fmt string, args ...interface{}) {
	logrus.Errorf(fmt, args...)
}

func Infof(fmt string, args ...interface{}) {
	logrus.Infof(fmt, args...)
}

func Fatalf(fmt string, args ...interface{}) {
	logrus.Fatalf(fmt, args...)
}
