package main

import (
	"go-watcher/daemon/config"
	"go-watcher/daemon/server"
	"log"
)

func main() {
	daemon := server.NewServer(config.NewServerConfig())
	if err := daemon.Start(); err != nil {
		log.Fatal(err)
	}
}
