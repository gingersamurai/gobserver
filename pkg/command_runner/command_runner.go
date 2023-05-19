package command_runner

import (
	"bytes"
	"os/exec"
)

type CommandRunner struct{}

func NewCommandRunner() CommandRunner {
	return CommandRunner{}
}

func (cr CommandRunner) RunCommand(command, dir string) (string, string, int, error) {
	cmd := exec.Command("bash", "-c", command)
	cmd.Dir = dir
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Start()
	if err != nil {
		return "", "", 0, err
	}
	err = cmd.Wait()
	if err != nil {
		return "", "", 0, err
	}
	return stdout.String(), stderr.String(), cmd.ProcessState.ExitCode(), nil // вернуть код выхода
}
