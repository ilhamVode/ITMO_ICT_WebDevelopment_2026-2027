package main

import (
	"fmt"
	"net"
)

func main() {
	fmt.Println("\n-- CLIENT --")

	addr, err := net.ResolveTCPAddr("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при разрешении адреса:", err)
		return
	}

	conn, err := net.DialTCP("tcp", nil, addr)
	if err != nil {
		fmt.Println("Ошибка при подключении:", err)
		return
	}
	defer conn.Close()

	buffer := make([]byte, 1024)

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Соединение закрыто")
			return
		}

		fmt.Print(string(buffer[:n]))

		var operation string

		fmt.Scan(&operation)

		if operation == "Выход" {
			conn.Write([]byte(operation))
			return
		}

		var message string

		switch operation {
		case "1", "4":
			var a, b float64
			fmt.Scan(&a, &b)

			message = fmt.Sprintf("%s %v %v", operation, a, b)

		case "2", "3":
			var a, b, c float64
			fmt.Scan(&a, &b, &c)

			message = fmt.Sprintf("%s %v %v %v", operation, a, b, c)

		default:
			fmt.Println("Неизвестная операция")
			continue
		}

		_, err = conn.Write([]byte(message))
		if err != nil {
			fmt.Println("Ошибка отправки:", err)
			return
		}

		if message == "Выход" {
			return
		}
	}
}
