package cmd

import (
	"bufio"
	"fmt"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"os"
	"strings"
)

func ParseCommand(command string) (string, string, error) {
	input := strings.Split(command[:len(command)-1], " ")
	if len(input) < 2 {
		return "", "", errors.New("input error")
	}
	c := input[0]
	p := input[1]
	return c, p, nil
}

func TerminalUI() {
	// todo 复用一个连接
	for {
		fmt.Print("watcherctl> ")
		inputReader := bufio.NewReader(os.Stdin)
		input, _ := inputReader.ReadString('\n')
		if input == "q\n" || input == "exit\n" {
			fmt.Println("bye!")
			break
		}
		requestType, process, _ := ParseCommand(input)
		tcpRequest := NewTcpRequest(requestType, process)
		if err := TcpClient(tcpRequest); err != nil {
			log.Fatal(err)
		}
	}
}
