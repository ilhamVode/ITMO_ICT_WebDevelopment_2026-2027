package teacher

import (
	"laboratory_work_3/internal/subject"
	"laboratory_work_3/internal/teaching"
)

type Teacher struct {
	ID               int64   `gorm:"primaryKey;column:id" json:"id"`
	FirstName        string  `gorm:"column:first_name" json:"first_name"`
	LastName         string  `gorm:"column:last_name" json:"last_name"`
	MiddleName       *string `gorm:"column:middle_name" json:"middle_name"`
	FixedClassroomID *int64  `gorm:"column:fixed_classroom_id" json:"fixed_classroom_id"`
	IsActive         bool    `gorm:"column:is_active" json:"is_active"`

	Assignments []teaching.TeachingAssignment `gorm:"foreignKey:TeacherID;references:ID" json:"assignments,omitempty"`
	Subjects    []subject.Subject             `gorm:"-" json:"subjects,omitempty"`
}
