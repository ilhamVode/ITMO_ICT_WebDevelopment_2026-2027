package main

import (
	"fmt"
	"math"
	"net"
	"strconv"
	"strings"
)

func main() {
	fmt.Println("\n-- SERVER --")

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

	clientConn, err := listener.AcceptTCP()
	if err != nil {
		fmt.Println("Ошибка при принятии соединения:", err)
		return
	}
	defer clientConn.Close()

	buf := make([]byte, 1024)

	clientConn.Write([]byte(
		"Добро пожаловать на сервер!\n\n" +
			"Выберите действие:\n" +
			"1. Вычисление гипотенузы\n" +
			"2. Решение квадратного уравнения\n" +
			"3. Вычисление площади трапеции\n" +
			"4. Вычисление площади параллелограмма\n" +
			"Введите номер действия и параметры через пробел.\n" +
			"Пример: 1 3 4\n" +
			"Для выхода введите 'Выход'.\n",
	))
	for {
		n, err := clientConn.Read(buf)
		if err != nil {
			fmt.Println("Ошибка при чтении данных от клиента:", err)
			return
		}

		message := strings.TrimSpace(string(buf[:n]))

		if message == "Выход" {
			fmt.Println("Клиент завершил соединение.")
			clientConn.Write([]byte("Вы закрыли соединение."))
			break
		}

		parts := strings.Fields(message)

		if len(parts) < 3 {
			clientConn.Write([]byte("Ошибка: недостаточно параметров"))
			continue
		}

		outResult, outResult2 := programCalculation(parts[0], parts[1:]...)
		var Result string
		if outResult2 != 0 {
			Result = fmt.Sprintf("Результат: %v, %v", outResult, outResult2)
		} else {
			Result = fmt.Sprintf("Результат: %v", outResult)
		}

		_, err = clientConn.Write([]byte(Result))
		if err != nil {
			fmt.Println("Ошибка при отправке данных клиенту:", err)
			return
		}
	}
}

func Pythagoras(a, b float64) float64 {
	return math.Sqrt((a * a) + (b * b))
}

func quadraticEquation(a, b, c float64) (float64, float64) {
	discriminant := (b * b) - (4 * a * c)
	if discriminant < 0 {
		return 0, 0
	}
	sqrtDiscriminant := math.Sqrt(discriminant)
	x1 := (-b + sqrtDiscriminant) / (2 * a)
	x2 := (-b - sqrtDiscriminant) / (2 * a)
	return x1, x2
}

func areaOfTrapezoid(a, b, h float64) float64 {
	return ((a + b) / 2) * h
}

func areaOfParallelogram(base, height float64) float64 {
	return base * height
}

func programCalculation(choice string, params ...string) (float64, float64) {
	paramsFloat := make([]float64, len(params))
	for i, param := range params {
		var value float64
		value, err := strconv.ParseFloat(param, 64)
		if err != nil {
			return 0, 0
		}
		paramsFloat[i] = value
	}
	switch choice {
	case "1":
		if len(params) >= 2 {
			return Pythagoras(paramsFloat[0], paramsFloat[1]), 0
		}
	case "2":
		if len(params) >= 3 {
			return quadraticEquation(paramsFloat[0], paramsFloat[1], paramsFloat[2])
		}
	case "3":
		if len(params) >= 3 {
			return areaOfTrapezoid(paramsFloat[0], paramsFloat[1], paramsFloat[2]), 0
		}
	case "4":
		if len(params) >= 2 {
			return areaOfParallelogram(paramsFloat[0], paramsFloat[1]), 0
		}
	}
	return 0, 0
}
