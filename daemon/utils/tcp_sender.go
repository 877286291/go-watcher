package utils

import (
	"net"
)

func WriteMsgToConn(conn *net.TCPConn, msg string) error {
	_, err := conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	if err = conn.CloseWrite(); err != nil {
		return err
	}
	return nil
}
