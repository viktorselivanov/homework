package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	address := "localhost:4242"
	timeout := 10 * time.Second

	client := NewTelnetClient(address, timeout, os.Stdin, os.Stdout)

	if err := client.Connect(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	done := make(chan struct{}, 2)

	go func() {
		_ = client.Send()
		done <- struct{}{}
	}()

	go func() {
		_ = client.Receive()
		done <- struct{}{}
	}()

	// ждём завершения одной из сторон
	<-done

	_ = client.Close()
}
