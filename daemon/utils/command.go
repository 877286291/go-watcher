package utils

import (
	"os"
	"os/exec"
)

func GetFullPath(file string) (string, error) {
	path, err := exec.LookPath(file)
	if err != nil {
		return "", err
	}
	return path, nil
}
func KillProcess(pid int) error {
	process, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	if err = process.Kill(); err != nil {
		return err
	}
	return nil
}
