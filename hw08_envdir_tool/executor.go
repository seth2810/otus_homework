package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	if len(cmd) == 0 {
		return -1
	}

	name, args := cmd[0], cmd[1:]

	for k, e := range env {
		if e.NeedRemove {
			if err := os.Unsetenv(k); err != nil {
				fmt.Println("error while running command:", err)
				return 1
			}
		} else {
			if err := os.Setenv(k, e.Value); err != nil {
				fmt.Println("error while running command:", err)
				return 1
			}
		}
	}

	c := exec.Command(name, args...)

	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	c.Env = os.Environ()

	err := c.Run()

	if err == nil {
		return 0
	}

	var exitError *exec.ExitError

	if errors.As(err, &exitError) {
		return exitError.ExitCode()
	}

	fmt.Println("error while running command:", err)

	return 1
}
