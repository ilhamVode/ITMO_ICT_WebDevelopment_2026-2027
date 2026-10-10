package report

import (
	"laboratory_work_3/internal/group"
	"laboratory_work_3/internal/semester"
	"laboratory_work_3/internal/subject"
)

type StudentGrade struct {
	SubjectID   int64  `json:"subject_id"`
	SubjectName string `json:"subject_name"`
	Value       *int   `json:"value"`
}

type StudentResult struct {
	ID        int64          `json:"id"`
	FirstName string         `json:"first_name"`
	LastName  string         `json:"last_name"`
	IsActive  bool           `json:"is_active"`
	Grades    []StudentGrade `json:"grades" gorm:"-"`
	Average   *float64       `json:"average" gorm:"-"`
}

type SubjectResult struct {
	ID         int64    `json:"id"`
	Name       string   `json:"name"`
	GradeCount int      `json:"grade_count"`
	Average    *float64 `json:"average"`
}

type Statement struct {
	Group        group.StudyGroup  `json:"group"`
	Semester     semester.Semester `json:"semester"`
	Students     []StudentResult   `json:"students"`
	Subjects     []SubjectResult   `json:"subjects"`
	StudentCount int               `json:"student_count"`
	GradeCount   int               `json:"grade_count"`
	Average      *float64          `json:"average"`
}

// Данные для построения ведомости, без вложенных объектов.
type gradeRow struct {
	StudentID int64
	SubjectID int64
	Value     int
}

type statementData struct {
	Group    group.StudyGroup
	Semester semester.Semester
	Students []StudentResult
	Subjects []subject.Subject
	Grades   []gradeRow
}
