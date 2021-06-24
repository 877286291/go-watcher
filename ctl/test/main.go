package main

import (
	"errors"
	"github.com/rivo/tview"
	"strings"
)

func main() {
	app := tview.NewApplication()
	inputField := tview.NewInputField().
		SetLabel("watcherctl>: ")
	if err := app.SetRoot(inputField, false).Run(); err != nil {
		panic(err)
	}
}
func ParseCommand(command string) (string, string, error) {
	input := strings.Split(command[:len(command)-1], " ")
	if len(input) < 2 {
		return "", "", errors.New("input error")
	}
	c := input[0]
	p := input[1]
	return c, p, nil
}
