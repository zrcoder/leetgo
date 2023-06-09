package exec

import (
	"os"
	"os/exec"
)

func Run(dir, cmd string, args ...string) error {
	command := exec.Command(cmd, args...)
	if dir != "" {
		command.Dir = dir
	}
	command.Stderr = os.Stderr
	command.Stdout = os.Stdout
	command.Stdin = os.Stdin
	return command.Run()
}
