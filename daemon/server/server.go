package server

import (
	"fmt"
	"go-watcher/daemon/config"
	"net"
)

type Server struct {
	Config *config.DaemonServerConfig
}

func NewServer(config *config.DaemonServerConfig) *Server {
	return &Server{
		config,
	}
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
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			return err
		}
		go handle(conn)
	}
}


