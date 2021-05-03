package main

import (
	"bytes"
	"errors"
	"io"
	"io/ioutil"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var terminalZeros = string([]byte{'\x00'})

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)

	for _, info := range files {
		if info.IsDir() || strings.Contains(info.Name(), "=") {
			continue
		}

		name := info.Name()

		data, err := ioutil.ReadFile(path.Join(dir, name))
		if err != nil {
			return nil, err
		}

		if len(data) == 0 {
			env[name] = EnvValue{
				NeedRemove: true,
			}

			continue
		}

		buff := bytes.NewBuffer(data)

		line, err := buff.ReadString('\n')

		if err != nil && !errors.Is(err, io.EOF) {
			return nil, err
		}

		value := strings.ReplaceAll(strings.TrimRight(line, "\n\t "), terminalZeros, "\n")

		env[name] = EnvValue{
			Value: value,
		}
	}

	return env, nil
}
