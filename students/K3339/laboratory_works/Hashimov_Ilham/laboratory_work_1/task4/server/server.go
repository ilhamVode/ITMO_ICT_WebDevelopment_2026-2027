package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

var (
	clients []*net.TCPConn // Срез для хранения подключений клиентов)
	mutex   sync.Mutex     // Мьютекс для синхронизации доступа к срезу клиентов
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

	for {
		clientConn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("Ошибка при принятии соединения:", err)
			continue
		}

		mutex.Lock()                          // Блокируем доступ к срезу клиентов
		clients = append(clients, clientConn) // Добавляем подключение клиента в срез
		mutex.Unlock()                        // Разблокируем доступ к срезу клиентов

		go handleClient(clientConn) // Запускаем обработку клиента в отдельной горутине(многопоточность)
	}

}

func handleClient(clientConn *net.TCPConn) {
	defer clientConn.Close()
	defer removeClient(clientConn) // Удаляем клиента из среза при завершении работы(то есть при закрытии соединения)

	buffer := make([]byte, 1024)

	for {
		n, err := clientConn.Read(buffer)
		if err != nil {
			fmt.Println("Ошибка при чтении данных от клиента:", err)
			return
		}

		fmt.Println("Началась обработка сообщения от клиента:", clientConn.RemoteAddr())

		time.Sleep(10 * time.Second) // Имитируем задержку обработки сообщения для демонстрации многопоточности

		message := fmt.Sprintf("CLIENT %s: %s", clientConn.RemoteAddr(), string(buffer[:n]))
		broadcastMessage(message) // Отправляем сообщение всем подключенным клиентам

		fmt.Println("Завершена обработка сообщения от клиента:", clientConn.RemoteAddr())
	}

}

// Функция для отправки сообщения всем подключенным клиентам
func broadcastMessage(message string) {
	mutex.Lock()
	defer mutex.Unlock()
	for _, client := range clients {
		_, err := client.Write([]byte(message))
		if err != nil {
			fmt.Println("Ошибка при отправке сообщения клиенту:", err)
		}
	}
}

func removeClient(clientConn *net.TCPConn) {
	mutex.Lock()
	defer mutex.Unlock()
	for i, client := range clients {
		if client == clientConn {
			clients = append(clients[:i], clients[i+1:]...) // Удаляем клиента из среза
			break
		}
	}
}
