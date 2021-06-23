package server

import (
	"go-watcher/daemon/config"
	"gopkg.in/ini.v1"
	"path/filepath"
)

func loadProcesses() ([]Process, error) {
	files, err := filepath.Glob(config.NewServerConfig().IncludePath)
	if err != nil {
		return nil, err
	}
	var programList []Process
	for _, file := range files {
		configFile, err := ini.Load(file)
		if err != nil {
			continue
		}
		programName := configFile.Section("program").Key("name").String()
		homeDir := configFile.Section("program").Key("home").String()
		command := configFile.Section("program").Key("command").String()
		arguments := configFile.Section("program").Key("arguments").String()
		autostart := configFile.Section("program").Key("autostart").MustBool(true)
		autoRestart := configFile.Section("program").Key("auto_restart").MustBool(false)
		retries := configFile.Section("program").Key("retries").MustInt(3)
		startSecs := configFile.Section("program").Key("startsecs").MustInt(3)
		environments := configFile.Section("program").Key("environment").MustString("")
		programList = append(programList, Process{
			ProcessName:   programName,
			HomeDirectory: homeDir,
			Command:       command,
			Arguments:     arguments,
			AutoStart:     autostart,
			AutoRestart:   autoRestart,
			Retries:       retries,
			StartSecs:     startSecs,
			Environment:   filepath.SplitList(environments),
			Status:        "Fatal",
		})
	}
	return programList, nil
}
