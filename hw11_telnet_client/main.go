package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var timeout time.Duration

var (
	ErrEOF                    = errors.New("...EOF")
	ErrConnectionClosedByPeer = errors.New("...Connection was closed by peer")
	ErrTooFewArguments        = errors.New("too few arguments in program call")
)

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "client connection timeout")
}

func receiveRoutine(client TelnetClient, errCh chan<- error) {
	err := client.Receive()

	if err == nil {
		err = ErrConnectionClosedByPeer
	}

	errCh <- err
}

func sendRoutine(client TelnetClient, errCh chan<- error) {
	err := client.Send()

	if err == nil {
		err = ErrEOF
	}

	errCh <- err
}

func runClient(address string) error {
	ctx, cancelFn := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGINT)

	defer cancelFn()

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		return fmt.Errorf("client connection error: %w", err)
	}

	defer client.Close()

	fmt.Fprintln(os.Stderr, "...Connected to", address)

	errCh := make(chan error, 1)

	go sendRoutine(client, errCh)
	go receiveRoutine(client, errCh)

	select {
	case <-ctx.Done():
		fmt.Println("client stopped after receiving SIGINT signal")
		return nil
	case err := <-errCh:
		return err
	}
}

func main() {
	flag.Parse()

	if flag.NArg() < 2 {
		fmt.Println(ErrTooFewArguments)
		os.Exit(1)
	}

	address := net.JoinHostPort(flag.Arg(0), flag.Arg(1))

	if err := runClient(address); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	os.Exit(0)
}
