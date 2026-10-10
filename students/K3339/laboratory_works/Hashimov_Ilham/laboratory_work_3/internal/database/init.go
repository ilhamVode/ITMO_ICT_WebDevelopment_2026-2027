package database

import (
	"fmt"
	"os"

	"gorm.io/gorm"
)

func InitSchema(db *gorm.DB) error {
	return db.Transaction(func(tx *gorm.DB) error {
		for _, filename := range []string{"migrations/init.sql", "migrations/business_rules.sql"} {
			script, err := os.ReadFile(filename)
			if err != nil {
				return fmt.Errorf("ошибка чтения %s: %w", filename, err)
			}
			if err := tx.Exec(string(script)).Error; err != nil {
				return fmt.Errorf("ошибка применения %s: %w", filename, err)
			}
		}
		return nil
	})
}
