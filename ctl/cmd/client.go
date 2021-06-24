package cmd

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net"
)

func TcpClient(tcpRequest *TcpRequest) error {
	tcpReqJsonStr, err := json.Marshal(tcpRequest)
	if err != nil {
		return err
	}
	tcpAddr, err := net.ResolveTCPAddr("tcp4", "127.0.0.1:2021")
	if err != nil {
		return err
	}
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	if err != nil {
		return err
	}
	defer conn.Close()                 // 关闭连接
	_, err = conn.Write(tcpReqJsonStr) // 发送数据
	if err != nil {
		return err
	}
	if err = conn.CloseWrite(); err != nil {
		return err
	}
	reader := bufio.NewReader(conn)
	for {
		slice, err := reader.ReadSlice('\n')
		if err != nil {
			break
		}
		fmt.Printf("%s", slice)
	}
	return nil
}
