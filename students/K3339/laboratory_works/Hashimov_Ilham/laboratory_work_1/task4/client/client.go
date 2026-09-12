package main

import (
	"bufio"
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

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Ошибка при подключении к серверу:", err)
		return
	}
	defer conn.Close()

	go receiveMessages(conn) // Запускаем горутину для получения сообщений от сервера(многопоточность)
	sendMessages(conn)       // здесь go не нужен, чтобы main не завершился раньше, чем завершится sendMessages

}

func receiveMessages(conn *net.TCPConn) {
	reader := bufio.NewReader(conn)
	for {
		message, err := reader.ReadString('\n') // Читаем сообщение до символа новой строки
		if err != nil {
			fmt.Println("Ошибка при чтении данных от сервера:", err)
			return
		}
		fmt.Print("Получено сообщение от сервера:", message)
	}
}

func sendMessages(conn *net.TCPConn) {
	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("Введите сообщение для отправки на сервер: ")

		scanner.Scan()
		message := scanner.Text()

		_, err := conn.Write([]byte(message + "\n")) // Добавляем символ новой строки для разделения сообщений
		if err != nil {
			fmt.Println("Ошибка при отправке сообщения на сервер:", err)
			return
		}
	}
}
