package teaching

import "laboratory_work_3/internal/subject"

type TeachingAssignment struct {
	ID         int64 `gorm:"primaryKey;column:id" json:"id"`
	TeacherID  int64 `gorm:"column:teacher_id" json:"teacher_id"`
	SubjectID  int64 `gorm:"column:subject_id" json:"subject_id"`
	GroupID    int64 `gorm:"column:group_id" json:"group_id"`
	SemesterID int64 `gorm:"column:semester_id" json:"semester_id"`

	Subject *subject.Subject `gorm:"foreignKey:SubjectID;references:ID" json:"subject,omitempty"`
}
