package main

import (
	"errors"
	"io"
	"net"
	"time"
)

var (
	ErrNoSuchHost     = errors.New("no such host")
	ErrConnectTimeout = errors.New("connect timeout")
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &client{address, timeout, in, out, nil}
}

type client struct {
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
	conn    io.ReadWriteCloser
}

func (c *client) Connect() (err error) {
	c.conn, err = net.DialTimeout("tcp", c.address, c.timeout)

	var dnsErr *net.DNSError

	if errors.As(err, &dnsErr) && dnsErr.IsNotFound {
		err = ErrNoSuchHost
	}

	var netErr net.Error

	if errors.As(err, &netErr) && netErr.Timeout() {
		err = ErrConnectTimeout
	}

	return
}

func (c *client) Send() error {
	_, err := io.Copy(c.conn, c.in)

	return err
}

func (c *client) Receive() error {
	_, err := io.Copy(c.out, c.conn)

	return err
}

func (c *client) Close() error {
	return c.conn.Close()
}
