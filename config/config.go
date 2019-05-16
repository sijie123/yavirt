package config

import (
	"path/filepath"

	"github.com/juju/errors"

	"github.com/projecteru2/yavirt/log"
)

const DefaultTemplate = `
env = "dev"
prof_http_port = 9999
bind_addr = "0.0.0.0:9696"

virt_dir = "/tmp/virt"
virt_bridge = "virbr0"

log_level = "warn"

mysql_timeout = "5s"
mysql_db = "test"
`

var Conf = newDefault()

type Config struct {
	Env          string
	ProfHttpPort int    `toml:"prof_http_port"`
	BindAddr     string `toml:"bind_addr"`

	VirtDir      string `toml:"virt_dir"`
	VirtFlockDir string
	VirtTmplDir  string
	VirtSockDir  string
	VirtBridge   string `toml:"virt_bridge"`

	LogLevel string `toml:"log_level"`
	LogFile  string `toml:"log_file"`

	MysqlUser     string   `toml:"mysql_user"`
	MysqlPassword string   `toml:"mysql_password"`
	MysqlAddr     string   `toml:"mysql_addr"`
	MysqlDB       string   `toml:"mysql_db"`
	MysqlTimeout  Duration `toml:"mysql_timeout"`
}

func newDefault() Config {
	var conf Config
	if err := Decode(DefaultTemplate, &conf); err != nil {
		log.Warnf(errors.ErrorStack(err))
		panic(err)
	}

	conf.loadVirtDirs()

	return conf
}

func (c *Config) Load(files []string) error {
	for _, path := range files {
		if err := c.load(path); err != nil {
			return errors.Trace(err)
		}
	}
	return nil
}

func (c *Config) load(file string) error {
	if err := DecodeFile(file, c); err != nil {
		return errors.Trace(err)
	}

	c.loadVirtDirs()

	return nil
}

func (c *Config) loadVirtDirs() {
	c.VirtFlockDir = filepath.Join(c.VirtDir, "flock")
	c.VirtTmplDir = filepath.Join(c.VirtDir, "template")
	c.VirtSockDir = filepath.Join(c.VirtDir, "sock")
}
