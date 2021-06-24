package server

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
	"go-watcher/daemon/utils"
	"io/ioutil"
	"net"
	"sort"
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
					_, err := conn.Write([]byte(fmt.Sprintf("process %s start %s\n", process.ProcessName, startStatus)))
					if err != nil {
						log.Error(err)
					}
					return
				}
				log.Infof("%s process already started\n", process.ProcessName)
				_, err := conn.Write([]byte(fmt.Sprintf("process %s already started\n", process.ProcessName)))
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
			if err := utils.WriteMsgToConn(conn, fmt.Sprintf("process %s already started\n", p.ProcessName)); err != nil {
				log.Error(err)
			}
		case "Stopped", "Fatal":
			startStatus := "success"
			if err := p.startProcess(); err != nil {
				startStatus = "failed"
			}
			if err := utils.WriteMsgToConn(conn, fmt.Sprintf("process %s start %s\n", p.ProcessName, startStatus)); err != nil {
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
				log.Infof("process %s start %s\n", r.ProcessName, startStatus)
				_, err := conn.Write([]byte(fmt.Sprintf("process %s start %s\n", r.ProcessName, startStatus)))
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
				if process.Pid == 0 {
					_, err := conn.Write([]byte(fmt.Sprintf("process %s is not running\n", process.ProcessName)))
					if err != nil {
						log.Errorln(err)
					}
				}
				killStatus := "success"
				if err := process.stopProcess(); err != nil {
					killStatus = "failed"
					log.Error(err)
				}
				log.Infof("%s process stop %s\n", process.ProcessName, killStatus)
				_, err := conn.Write([]byte(fmt.Sprintf("process %s stop %s\n", process.ProcessName, killStatus)))
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
				if err := utils.WriteMsgToConn(conn, fmt.Sprintf("process %s not running\n", r.ProcessName)); err != nil {
					log.Errorln(err)
				}
				return
			}
			if err := utils.WriteMsgToConn(conn, fmt.Sprintf("process %s stop failed\n", r.ProcessName)); err != nil {
				log.Errorln(err)
			}
			return
		}
		if err := utils.WriteMsgToConn(conn, fmt.Sprintf("process %s stop success\n", r.ProcessName)); err != nil {
			log.Errorln(err)
		}
		return
	}
	if err := utils.WriteMsgToConn(conn, fmt.Sprintf("no such process %s\n", r.ProcessName)); err != nil {
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
				_, err := conn.Write([]byte(fmt.Sprintf("process %s stop %s\n", r.ProcessName, status)))
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
		if err := utils.WriteMsgToConn(conn, fmt.Sprintf("process %s restart %s", r.ProcessName, status)); err != nil {
			log.Error(err)
		}
	}
}
func (r TcpRequest) getProcessStatus(conn *net.TCPConn) {
	log.Info("getProcessStatus request")
	if r.ProcessName == "all" {
		keys := make([]string, 0, len(ProcessRegister))
		for k := range ProcessRegister {
			keys = append(keys, k)
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		for _, key := range keys {
			status := ProcessRegister[key].getProcessStatus()
			_, err := conn.Write([]byte(status))
			if err != nil {
				log.Errorln(err)
			}
		}
		if err := conn.CloseWrite(); err != nil {
			log.Errorln(err)
		}
		return
	}
	if processStatus, ok := ProcessRegister[r.ProcessName]; ok {
		status := processStatus.getProcessStatus()
		if err := utils.WriteMsgToConn(conn, status); err != nil {
			log.Errorln(err)
		}
		return
	}
	if err := utils.WriteMsgToConn(conn, fmt.Sprintf("no such process %s\n", r.ProcessName)); err != nil {
		log.Errorln(err)
	}
}
