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

	// Пытаемся подключиться к серверу.
	if err := client.Connect(); err != nil {
		// Если подключение не удалось, выводим ошибку в stderr и завершаем программу
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	// Конкурентно запускаем горутину для отправки данных с stdin в сокет.
	go func() {
		if err := client.Send(); err != nil {
			// Если произошла ошибка при отправке данных, выводим её в stderr
			fmt.Fprintln(os.Stderr, err)
		}
		client.Close()
	}()

	if err := client.Receive(); err != nil {
		fmt.Fprintln(os.Stderr, "Receive finished:", err) // Обычно Receive завершится, когда сервер закроет соединение
	}

	client.Close()
}
