package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("-- CLIENT --")
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при разрешении адреса:", err)
		return
	}

	// Создание UDP сокета для клиента
	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		fmt.Println("Ошибка при создании UDP сокета:", err)
		return
	}
	defer conn.Close()

	_, err = conn.Write([]byte("Hello, server!"))
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении сообщения:", err)
		return
	}
	fmt.Println("Получено сообщение от сервера:", string(buffer[:n]))
}
