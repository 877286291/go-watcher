package server

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-watcher/daemon/config"
	"net"
	"os"
	"path/filepath"
)

type Server struct {
	Config *config.DaemonServerConfig
}

func NewServer(config *config.DaemonServerConfig) *Server {
	return &Server{
		config,
	}
}
func InitProcess() error {
	processes, err := loadProcesses()
	if err != nil {
		return err
	}
	for _, process := range processes {
		process := process
		ProcessRegister[process.ProcessName] = &process
		if process.AutoStart {
			process.Environment = append(process.Environment, filepath.SplitList(os.Getenv("PATH"))...)
			go func() {
				if err := process.startProcess(); err != nil {
					log.Errorf("process %s start failed", process.ProcessName)
					return
				}
				log.Infof("process %s start success", process.ProcessName)
			}()
		}
	}
	return nil
}
func (s *Server) Start() error {
	addr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf(":%d", s.Config.ServerPort))
	if err != nil {
		return err
	}
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}
	log.Info("Listen on port:", addr.Port)
	if err = InitProcess(); err != nil {
		return err
	}
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		go handle(conn)
	}
}
