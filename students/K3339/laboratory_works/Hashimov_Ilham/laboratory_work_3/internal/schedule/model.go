package schedule

import (
	"laboratory_work_3/internal/classroom"
	"laboratory_work_3/internal/teaching"
)

type Schedule struct {
	ID           int64 `gorm:"primaryKey;column:id" json:"id"`
	AssignmentID int64 `gorm:"column:assignment_id" json:"assignment_id"`
	ClassroomID  int64 `gorm:"column:classroom_id" json:"classroom_id"`
	Weekday      int   `gorm:"column:weekday" json:"weekday"`
	LessonNumber int   `gorm:"column:lesson_number" json:"lesson_number"`

	Assignment *teaching.TeachingAssignment `gorm:"foreignKey:AssignmentID;references:ID" json:"assignment,omitempty"`
	Classroom  *classroom.Classroom         `gorm:"foreignKey:ClassroomID;references:ID" json:"classroom,omitempty"`
}
