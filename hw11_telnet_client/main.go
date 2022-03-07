package main

import (
	"context"
	"flag"
	"fmt"
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
	args := flag.Args()
	if argLen := len(args); argLen != 2 {
		log.Fatalf("Invalid number of arguments (%d)\n", argLen)
	}

	address := net.JoinHostPort(args[0], args[1])

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)
	if err := client.Connect(); err != nil {
		log.Fatalf("Can't connect to %s\n", address)
	}
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
				_, _ = fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
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
