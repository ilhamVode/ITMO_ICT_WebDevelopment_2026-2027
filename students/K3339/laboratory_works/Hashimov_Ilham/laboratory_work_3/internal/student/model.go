package student

import "time"

type Student struct {
	ID        int64      `gorm:"primaryKey;column:id" json:"id"`
	FirstName string     `gorm:"column:first_name" json:"first_name"`
	LastName  string     `gorm:"column:last_name" json:"last_name"`
	BirthDate *time.Time `gorm:"column:birth_date;type:date" json:"birth_date"`
	GroupID   int64      `gorm:"column:group_id" json:"group_id"`
	IsActive  bool       `gorm:"column:is_active" json:"is_active"`
}
