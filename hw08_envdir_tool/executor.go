package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return 1
	}
	commandName := cmd[0]

	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	command := exec.Command(commandName, args...)

	for name, envValue := range env {
		err := os.Unsetenv(name)
		if err != nil {
			return 1
		}
		if envValue.NeedRemove {
			continue
		}
		err = os.Setenv(name, envValue.Value)
		if err != nil {
			return 1
		}
	}

	command.Env = os.Environ()
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr

	if err := command.Run(); err != nil {
		var ee *exec.ExitError
		if errors.As(err, &ee) {
			returnCode = ee.ExitCode()
		}
	}

	return
}
