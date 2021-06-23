package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-watcher/daemon/utils"
	"io/ioutil"
	"net"
	"sync"
)

type TcpRequest struct {
	RequestType string `json:"request_type"`           //请求类型 status start stop restart reload
	ProcessName string `json:"process_name,omitempty"` // 进程名称
}

func NewTcpRequest(requestType string, processName string) *TcpRequest {
	return &TcpRequest{
		RequestType: requestType,
		ProcessName: processName,
	}
}

func handle(conn *net.TCPConn) {
	defer conn.Close()
	request, err := ioutil.ReadAll(conn)
	if err != nil {
		log.Fatalln(err)
	}
	var tcpRequest TcpRequest
	if err := json.Unmarshal(request, &tcpRequest); err != nil {
		log.Fatalln(err)
	}
	switch tcpRequest.RequestType {
	case "status":
		tcpRequest.getProcessStatus(conn)
	case "start":
		tcpRequest.startProcess(conn)
	case "stop":
		tcpRequest.stopProcess(conn)
	case "restart":
		tcpRequest.restartProcess(conn)
	}
	return
}

func (r TcpRequest) startProcess(conn *net.TCPConn) {
	if r.ProcessName == "all" {
		var wg sync.WaitGroup
		for _, process := range ProcessRegister {
			wg.Add(1)
			go func(process *Process) {
				defer wg.Done()
				if process.Status != "Running" {
					startStatus := "success"
					if err := process.startProcess(); err != nil {
						startStatus = "failed"
					}
					_, err := conn.Write([]byte(fmt.Sprintf("%s process start %s\n", process.ProcessName, startStatus)))
					if err != nil {
						log.Error(err)
					}
					return
				}
				log.Infof("%s process already started\n", process.ProcessName)
				_, err := conn.Write([]byte(fmt.Sprintf("%s process already started\n", process.ProcessName)))
				if err != nil {
					log.Error(err)
				}
			}(process)
		}
		wg.Wait()
		if err := conn.CloseWrite(); err != nil {
			log.Errorln(err)
		}
		return
	}
	if p, ok := ProcessRegister[r.ProcessName]; ok {
		switch p.Status {
		case "Running":
			log.Infof("%s already started", p.ProcessName)
			if err := utils.WriteMsgToConn(conn, fmt.Sprintf("%s already started\n", p.ProcessName)); err != nil {
				log.Error(err)
			}
		case "Stopped", "Fatal":
			startStatus := "success"
			if err := p.startProcess(); err != nil {
				startStatus = "failed"
			}
			if err := utils.WriteMsgToConn(conn, fmt.Sprintf("%s process start %s\n", p.ProcessName, startStatus)); err != nil {
				log.Errorln(err)
			}
		}
		return
	}
	processes, err := loadProcesses()
	if err != nil {
		log.Error(err)
	}
	var wg sync.WaitGroup
	for _, process := range processes {
		wg.Add(1)
		go func(process Process) {
			defer wg.Done()
			if process.ProcessName == r.ProcessName {
				startStatus := "success"
				if err := process.startProcess(); err != nil {
					startStatus = "failed"
				}
				log.Infof("%s process start %s\n", r.ProcessName, startStatus)
				_, err := conn.Write([]byte(fmt.Sprintf("%s process start %s\n", r.ProcessName, startStatus)))
				if err != nil {
					log.Errorln(err)
				}
			}
		}(process)
	}
	wg.Wait()
	if err = conn.CloseWrite(); err != nil {
		log.Errorln(err)
	}
}
func (r TcpRequest) stopProcess(conn *net.TCPConn) {
	if r.ProcessName == "all" {
		var wg sync.WaitGroup
		for _, process := range ProcessRegister {
			wg.Add(1)
			go func(process *Process) {
				defer wg.Done()
				killStatus := "success"
				if err := process.stopProcess(); err != nil {
					killStatus = "failed"
					log.Error(err)
				}
				log.Infof("%s process stop %s\n", process.ProcessName, killStatus)
				_, err := conn.Write([]byte(fmt.Sprintf("%s process stop %s\n", process.ProcessName, killStatus)))
				if err != nil {
					log.Errorln(err)
				}
			}(process)
		}
		wg.Wait()
		if err := conn.CloseWrite(); err != nil {
			log.Errorln(err)
		}
		return
	}
	if p, ok := ProcessRegister[r.ProcessName]; ok {
		if err := p.stopProcess(); err != nil {
			if p.Status == "Stopped" {
				if err := utils.WriteMsgToConn(conn, fmt.Sprintf("%s process not running\n", r.ProcessName)); err != nil {
					log.Errorln(err)
				}
				return
			}
			if err := utils.WriteMsgToConn(conn, fmt.Sprintf("%s process stop failed\n", r.ProcessName)); err != nil {
				log.Errorln(err)
			}
			return
		}
		if err := utils.WriteMsgToConn(conn, fmt.Sprintf("%s process stop success\n", r.ProcessName)); err != nil {
			log.Errorln(err)
		}
		return
	}
	if err := utils.WriteMsgToConn(conn, fmt.Sprintf("no such process:%s\n", r.ProcessName)); err != nil {
		log.Errorln(err)
	}
}
func (r TcpRequest) restartProcess(conn *net.TCPConn) {
	if r.ProcessName == "all" {
		var wg sync.WaitGroup
		for _, process := range ProcessRegister {
			wg.Add(1)
			go func(process *Process) {
				defer wg.Done()
				status := "success"
				if err := process.restartProcess(); err != nil {
					status = "failed"
				}
				_, err := conn.Write([]byte(fmt.Sprintf("%s process stop %s\n", r.ProcessName, status)))
				if err != nil {
					log.Errorln(err)
				}
			}(process)
		}
		wg.Wait()
		if err := conn.CloseWrite(); err != nil {
			log.Error(err)
		}
		return
	}
	if p, ok := ProcessRegister[r.ProcessName]; ok {
		status := "success"
		if err := p.restartProcess(); err != nil {
			status = "failed"
			log.Error(err)
		}
		if err := utils.WriteMsgToConn(conn, fmt.Sprintf("%s process restart %s", r.ProcessName, status)); err != nil {
			log.Error(err)
		}
	}
}
func (r TcpRequest) getProcessStatus(conn *net.TCPConn) {

}
