package main

import (
	"bufio"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	fileInfos, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment)
	for _, fileInfo := range fileInfos {
		if strings.Contains(fileInfo.Name(), "=") {
			continue
		}

		if fileInfo.Size() == 0 {
			env[fileInfo.Name()] = EnvValue{"", true}
			continue
		}

		filePath := dir + "/" + fileInfo.Name()
		envValue, err := getEnvValue(filePath)
		if err != nil {
			return nil, err
		}

		env[fileInfo.Name()] = envValue
	}

	return env, nil
}

func getEnvValue(filePath string) (EnvValue, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return EnvValue{}, err
	}
	defer file.Close()

	fileScanner := bufio.NewScanner(file)
	if !fileScanner.Scan() {
		return EnvValue{}, fileScanner.Err()
	}

	return EnvValue{clean(fileScanner.Text()), false}, nil
}

func clean(content string) string {
	content = strings.TrimRight(content, " \t")
	return string(bytes.ReplaceAll([]byte(content), []byte{0x00}, []byte("\n")))
}
