package main

import (
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatal("Invalid number of arguments. There should be more than two of them")
	}

	envDir := os.Args[1]
	env, err := ReadDir(envDir)
	if err != nil {
		log.Fatal(err)
	}
	command := os.Args[2:]
	os.Exit(RunCmd(command, env))
}
