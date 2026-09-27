package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres создаёт подключение GORM к PostgreSQL
func NewPostgres(dsn string) (*gorm.DB, error) {
	// Создаём PostgreSQL-драйвер с параметрами подключения
	// т.е. дает понять то, что мы будем работать именно с postgres
	// так же задаются параметры подключения через dsn
	dialector := postgres.Open(dsn)

	// Создаём объект GORM для работы с PostgreSQL
	// тут уже dialector дает понять, что мы используем postgres
	db, err := gorm.Open(dialector, &gorm.Config{})

	if err != nil {
		return nil, err
	}

	// метод DB() позволяет достать *sql.DB из *gorm.DB
	// опуститься на низкий уровень нужно т.к. Ping() принадлежит *sql.DB
	sqlDB, err := db.DB()

	if err != nil {
		return nil, err
	}

	// Проверяем, что PostgreSQL действительно доступен
	// Ping() позволяет проверить доступна ли база данных
	err = sqlDB.Ping()

	if err != nil {
		return nil, err
	}

	return db, nil
}
