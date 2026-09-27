package main

import (
	"log"
	"net"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"laboratory_work_2/internal/app"
	"laboratory_work_2/internal/database"
)

func main() {
	log.Println("Загрузка конфигурации...")

	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("Отсутствует переменная DATABASE_DSN")
	}

	log.Println("Подключение к PostgreSQL...")

	db, err := database.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Подключение к PostgreSQL прошло успешно!")

	log.Println("Создание репозиториев...")
	userRepo := app.NewPostgresUserRepository(db)
	homeworkRepo := app.NewPostgresHomeworkRepository(db)
	submissionRepo := app.NewPostgresSubmissionRepository(db)
	log.Println("Репозитории успешно созданы")

	sessions := app.NewSessionManager()
	auth := app.NewAuth(userRepo, sessions)

	homeworkView := app.NewHomeworkView(homeworkRepo, auth)
	adminView := app.NewAdminView(userRepo, homeworkRepo, submissionRepo)
	userView := app.NewUserView(userRepo)
	router := app.NewRouter(homeworkView, adminView, userView, auth)

	addr := "127.0.0.1:8080"

	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("Сервер запущен на http://%s", listener.Addr().String())

	err = http.Serve(listener, router)
	if err != nil {
		log.Fatal(err)
	}

}
