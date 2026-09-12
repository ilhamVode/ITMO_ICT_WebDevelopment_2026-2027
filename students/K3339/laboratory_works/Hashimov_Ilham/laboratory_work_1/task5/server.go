package main

import (
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"
)

func main() {
	listener, err := net.Listen("tcp", "localhost:8080")
	if err != nil {
		fmt.Println("Ошибка при прослушивании порта:", err)
		return
	}
	defer listener.Close()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Ошибка при приеме соединения:", err)
			return
		}

		handleClient(conn)
	}
}

func handleClient(conn net.Conn) {

	defer conn.Close()

	buf := make([]byte, 1024)

	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Ошибка при чтении данных:", err)
		return
	}

	request := string(buf[:n])

	parts := strings.Fields(request)

	// Проверяем, что запрос содержит хотя бы три части: метод, путь и версию протокола
	if len(parts) < 3 {
		fmt.Println("Некорректный запрос")
		return
	}

	method := parts[0]
	path := parts[1]

	method = strings.ToUpper(method)
	if method != "GET" && method != "POST" {
		fmt.Println("Некорректный метод:", method)
		return
	} else if method == "GET" {
		handleGet(conn, path)
	} else if method == "POST" {
		handlePost(conn, path, request)
	}
}

func handleGet(conn net.Conn, path string) {
	html, err := os.ReadFile("index.html")
	if err != nil {
		fmt.Println("Ошибка при чтении файла:", err)
		return
	}

	gradesHTML := ""
	for _, grade := range grades {
		gradesHTML += "<li>" + grade.Subject + " - " + grade.Grade + "</li>"
	}

	// Заменяем плейсхолдер {{GRADES}} на сгенерированный HTML-код с оценками
	htmlString := strings.Replace(
		string(html),
		"{{GRADES}}",
		gradesHTML,
		1,
	)

	response := "HTTP/1.1 200 OK\r\n" +
		"Content-Type: text/html; charset=utf-8\r\n" +
		"Content-Length: " + strconv.Itoa(len(htmlString)) + "\r\n" +
		"\r\n" +
		htmlString

	_, err = conn.Write([]byte(response))
	if err != nil {
		fmt.Println("Ошибка при отправке ответа:", err)
	}
}

type Grade struct {
	Subject string
	Grade   string
}

var grades []Grade

func handlePost(conn net.Conn, path string, request string) {
	parts := strings.Split(request, "\r\n\r\n")

	// не содержит тело
	if len(parts) < 2 {
		fmt.Println("Некорректный POST-запрос")
		return
	}

	body := parts[1]

	// Разделяем тело запроса на поля
	fields := strings.Split(body, "&")

	var subject, grade string

	for _, field := range fields {
		pair := strings.Split(field, "=")

		// Проверяем, что поле содержит ключ и значение
		if len(pair) != 2 {
			fmt.Println("Некорректное поле:", field)
			return
		}

		switch pair[0] {
		case "subject":
			subject = pair[1]
		case "grade":
			grade = pair[1]
		}
	}

	grades = append(grades, Grade{
		Subject: subject,
		Grade:   grade,
	})

	fmt.Println("Добавлена оценка:", subject, grade)

	// браузер сразу же перенаправит на главную страницу
	handleGet(conn, "/")
}
