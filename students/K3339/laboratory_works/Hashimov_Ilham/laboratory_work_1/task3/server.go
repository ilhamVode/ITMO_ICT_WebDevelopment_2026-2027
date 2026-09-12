package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при разрешении адреса:", err)
		return
	}

	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		fmt.Println("Ошибка при создании TCP-соединения:", err)
		return
	}
	defer listener.Close()

	fmt.Println("Сервер запущен. Ожидание клиента...")

	for {
		clientConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		handleClient(clientConn)
	}
}

func handleClient(conn *net.TCPConn) {
	defer conn.Close()

	buffer := make([]byte, 1024)

	n, err := conn.Read(buffer)
	if err != nil {
		fmt.Println("Ошибка при чтении данных от клиента:", err)
		return
	}

	fmt.Println("HTTP запрос:")
	fmt.Println(string(buffer[:n]))

	html, err := os.ReadFile("task3/index.html")
	if err != nil {
		fmt.Println("Ошибка при чтении файла index.html:", err)
		return
	}

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html\r\n" +
		"Content-Length: " + fmt.Sprint(len(html)) + "\r\n" +
		"\r\n" +
		string(html)

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа клиенту:", err)
		return
	}
}
