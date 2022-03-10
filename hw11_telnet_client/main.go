package main

import (
	"context"
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"sync"
	"time"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", time.Second+10, "Connection timeout")
}

func main() {
	flag.Parse()

	address := getAddress(flag.Args())
	client := *getTelnetClient(address)
	defer func() { _ = client.Close() }()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	go func() {
		<-ctx.Done()
		stop()
		os.Exit(1)
	}()

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for {
			err := client.Send()
			if err != nil {
				os.Exit(1)
			}
		}
	}()

	go func() {
		defer wg.Done()
		for {
			err := client.Receive()
			if err != nil {
				os.Exit(1)
			}
		}
	}()

	wg.Wait()
}

func getTelnetClient(address string) *TelnetClient {
	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Can't connect to %s\n", address)
	}

	return &client
}

func getAddress(args []string) string {
	if argLen := len(args); argLen != 2 {
		log.Fatalf("Invalid number of arguments (%d)\n", argLen)
	}

	return net.JoinHostPort(args[0], args[1])
}
