package main

import (
	"errors"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

// Реализация.
type telnetClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      io.ReadCloser
	out     io.Writer
	done    chan struct{}
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClient{
		address: address,
		timeout: timeout,
		in:      in,
		out:     out,
		done:    make(chan struct{}),
	}
}

// Connect подключается к серверу.
func (t *telnetClient) Connect() error {
	conn, err := net.DialTimeout("tcp", t.address, t.timeout)
	if err != nil {
		return fmt.Errorf("connection failed: %w", err)
	}
	t.conn = conn

	// выводим служебное сообщение в stderr, не в out
	fmt.Fprintf(os.Stderr, "...Connected to %s\n", t.address)
	return nil
}

// Close закрывает соединение.
func (t *telnetClient) Close() error {
	if t.conn != nil {
		err := t.conn.Close()
		close(t.done)
		return err
	}
	return nil
}

// Send читает из stdin и пишет в сокет.
func (t *telnetClient) Send() error {
	_, err := io.Copy(t.conn, t.in)
	if err != nil {
		if errors.Is(err, io.EOF) {
			fmt.Fprintln(t.out, "...EOF")
			return nil
		}
		return err
	}
	return nil
}

// Receive читает из сокета и пишет в stdout.
func (t *telnetClient) Receive() error {
	_, err := io.Copy(t.out, t.conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, "...Connection was closed by peer")
		return err
	}
	return nil
}
