package database

import (
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// NewPostgres создаёт подключение GORM к PostgreSQL
func NewPostgres(dsn string) (*gorm.DB, error) {
	// Создаём PostgreSQL-драйвер с параметрами подключения
	// PreferSimpleProtocol позволяет выполнять несколько SQL-команд
	// из init.sql одним запросом
	dialector := postgres.New(postgres.Config{
		DSN:                  dsn,
		PreferSimpleProtocol: true,
	})

	// Создаём объект GORM для работы с PostgreSQL
	db, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// DB() позволяет получить *sql.DB из *gorm.DB
	// Это нужно, потому что Ping() принадлежит *sql.DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, err
	}

	// Проверяем, что PostgreSQL действительно доступен
	if err := sqlDB.Ping(); err != nil {
		sqlDB.Close()
		return nil, err
	}

	return db, nil
}
