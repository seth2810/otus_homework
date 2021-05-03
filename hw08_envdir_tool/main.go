package main

import (
	"errors"
	"fmt"
	"os"
)

var ErrTooFewArguments = errors.New("too few arguments in program call")

func main() {
	if len(os.Args) <= 2 {
		fmt.Println(ErrTooFewArguments)
		os.Exit(1)
	}

	envDir, cmd := os.Args[1], os.Args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println("env dir reading error:", err)
		os.Exit(1)
	}

	if code := RunCmd(cmd, env); code != 0 {
		os.Exit(code)
	}
}
