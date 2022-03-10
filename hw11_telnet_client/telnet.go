package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

var ErrClosedChannel = errors.New("the channel is close")

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TClient{Address: address, Timeout: timeout, In: in, Out: out}
}

type TClient struct {
	Address string
	Timeout time.Duration
	In      io.ReadCloser
	Out     io.Writer
	Conn    net.Conn
}

func (t *TClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.Address, t.Timeout)
	if err == nil {
		t.Conn = conn
		t.log(fmt.Sprintf("...Connected to %s", t.Address))
	}

	return err
}

func (t *TClient) Close() error {
	return t.Conn.Close()
}

func (t *TClient) Send() error {
	err := t.call(t.In, t.Conn)
	if err != nil {
		if errors.Is(err, io.EOF) {
			t.log("...EOF")
			return nil
		}
		t.log("...Connection was closed by peer")
	}
	return err
}

func (t *TClient) Receive() error {
	err := t.call(t.Conn, t.Out)
	if err != nil {
		if errors.Is(err, io.EOF) {
			return nil
		}
		t.log("...Connection was closed by peer")
	}
	return err
}

func (t *TClient) call(r io.Reader, w io.Writer) error {
	ch, errch := makeChannel(r)
	for {
		select {
		case str, ok := <-ch:
			if !ok {
				return ErrClosedChannel
			}
			_, err := w.Write([]byte(str))
			if err != nil {
				return err
			}
		case errc := <-errch:
			return errc
		}
	}
}

func (t *TClient) log(msg string) {
	_, _ = fmt.Fprintln(os.Stderr, msg)
}

func makeChannel(r io.Reader) (<-chan string, <-chan error) {
	ch := make(chan string)
	errch := make(chan error)
	go func() {
		defer close(ch)
		reader := bufio.NewReader(r)
		for {
			text, err := reader.ReadString('\n')
			if err != nil {
				errch <- err
				return
			}
			ch <- text
		}
	}()
	return ch, errch
}
