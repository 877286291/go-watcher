package server

import (
	"errors"
	"fmt"
	"go-watcher/daemon/utils"
	"os"
	"time"
)

var ProcessRegister map[string]*Process

type Process struct {
	ProcessName   string    `json:"process_name"`
	HomeDirectory string    `json:"-"`
	Command       string    `json:"command"`
	Arguments     string    `json:"arguments"`
	AutoStart     bool      `json:"-"`
	AutoRestart   bool      `json:"-"`
	Retries       int       `json:"retries"`
	StartSecs     int       `json:"-"`
	Environment   []string  `json:"-"`
	Status        string    `json:"status"`
	Stop          bool      `json:"-"`
	Pid           int       `json:"pid,omitempty"`
	StartTime     time.Time `json:"start_time"`
	FatalCount    int       `json:"fatal_count"`
}

func init() {
	ProcessRegister = make(map[string]*Process)
}
func (p Process) wait4(process *os.Process) {
	p.Pid = process.Pid
	p.Status = "Running"
	tmp := p.FatalCount
	p.FatalCount = 0
	_, _ = process.Wait()
	p.FatalCount = tmp
	p.Pid = 0
	p.Status = "Fatal"
	// 手动停止进程
	if p.Stop {
		p.Status = "Stopped"
		p.FatalCount = 0
		p.Stop = false
		return
	}
	if p.Retries-1 > p.FatalCount {
		p.FatalCount++
		// 尝试拉起
		_ = p.startProcess()
	}
	// 失败重置
	p.FatalCount = 0
}
func (p Process) startProcess() error {
	procAttr := &os.ProcAttr{
		Dir:   p.HomeDirectory,
		Env:   p.Environment,
		Files: []*os.File{os.Stdin, os.Stdout, os.Stderr},
	}
	fullPath, err := utils.GetFullPath(p.Command)
	if err != nil || fullPath == "" {
		return errors.New(fmt.Sprintf("%s服务启动失败,请检查command字段是否正确\n", p.ProcessName))
	}
	args := []string{fullPath}
	if p.Arguments != "" {
		args = append(args, p.Arguments)
	}
	processStatus, err := os.StartProcess(fullPath, args, procAttr)
	if err != nil {
		return err
	}
	go p.wait4(processStatus)
	timer := time.After(time.Second * time.Duration(p.StartSecs))
	for {
		select {
		case <-timer:
			return nil
		}
	}
}
func (p Process) stopProcess() error {
	if p.Pid == 0 {
		return errors.New("process is not running")
	}
	p.Stop = true
	if err := utils.KillProcess(p.Pid); err != nil {
		p.Stop = false
		return err
	}
	return nil
}
func (p Process) restartProcess() error {
	if p.Pid != 0 {
		if err := p.stopProcess(); err != nil {
			return err
		}
	}
	if err := p.startProcess(); err != nil {
		return err
	}
	return nil
}
func (p Process) getProcessStatus() error {

	return nil
}
