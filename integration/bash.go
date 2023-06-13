package integration

import (
	"os/exec"
	"strings"
)

type BashProcess struct{}

func NewBashProcess() (*BashProcess, error) {
	return &BashProcess{}, nil
}

func (bp *BashProcess) Run(commands []string) (string, error) {
	command := strings.Join(commands, ";")

	cmd := exec.Command("bash", "-c", command) //nolint gosec

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(stdoutStderr), nil
}
