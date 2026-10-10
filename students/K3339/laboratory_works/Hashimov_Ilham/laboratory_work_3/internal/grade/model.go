package grade

import (
	"laboratory_work_3/internal/semester"
	"laboratory_work_3/internal/student"
	"laboratory_work_3/internal/subject"
)

type Grade struct {
	ID         int64 `gorm:"primaryKey;column:id" json:"id"`
	StudentID  int64 `gorm:"column:student_id" json:"student_id"`
	SubjectID  int64 `gorm:"column:subject_id" json:"subject_id"`
	SemesterID int64 `gorm:"column:semester_id" json:"semester_id"`
	Value      int   `gorm:"column:value" json:"value"`

	Student  *student.Student   `gorm:"foreignKey:StudentID;references:ID" json:"student,omitempty"`
	Subject  *subject.Subject   `gorm:"foreignKey:SubjectID;references:ID" json:"subject,omitempty"`
	Semester *semester.Semester `gorm:"foreignKey:SemesterID;references:ID" json:"semester,omitempty"`
}
