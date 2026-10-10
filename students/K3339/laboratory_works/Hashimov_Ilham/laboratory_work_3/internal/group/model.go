package group

import "laboratory_work_3/internal/student"

type StudyGroup struct {
	ID       int64             `gorm:"primaryKey;column:id" json:"id"`
	Name     string            `gorm:"column:name" json:"name"`
	Course   int               `gorm:"column:course" json:"course"`
	Students []student.Student `gorm:"foreignKey:GroupID;references:ID" json:"students,omitempty"`
}
