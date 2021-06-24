package config

import (
	log "github.com/sirupsen/logrus"
	"gopkg.in/ini.v1"
)

var cfg *ini.File

type DaemonServerConfig struct {
	ServerPort  int
	LogPath     string
	IncludePath string
}

func init() {
	ParseConfig("/Users/houyuji/go-watcher.ini")
}
func NewServerConfig() *DaemonServerConfig {
	return &DaemonServerConfig{
		ServerPort:  cfg.Section("go-watcher").Key("server_port").MustInt(2021),
		LogPath:     cfg.Section("go-watcher").Key("log_path").MustString("tmp"),
		IncludePath: cfg.Section("include").Key("files").String(),
	}
}
func DefaultServerConfig() *DaemonServerConfig {
	return &DaemonServerConfig{
		ServerPort: 2021,
		LogPath:    "tmp",
	}
}
func ParseConfig(cfgFile string) {
	var err error
	cfg, err = ini.Load(cfgFile)
	if err != nil {
		log.Fatalf("Fail to read file: %s", err)
	}
}
