package main

import (
	"fmt"
	"net"
)

func main() {

	fmt.Println("-- SERVER --")
	addr, err := net.ResolveUDPAddr("udp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при разрешении адреса:", err)
		return
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		fmt.Println("Ошибка при создании UDP сокета:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)
	n, clientAddr, err := conn.ReadFromUDP(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении сообщения:", err)
		return
	}
	fmt.Println("Получено сообщение от клиента:", string(buffer[:n]), "от", clientAddr)

	_, err = conn.WriteToUDP([]byte("Hello, client!"), clientAddr)
	if err != nil {
		fmt.Println("Ошибка при отправке сообщения:", err)
		return
	}
}
