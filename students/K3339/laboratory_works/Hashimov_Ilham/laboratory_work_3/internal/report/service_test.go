package report

import (
	"testing"

	"laboratory_work_3/internal/subject"
)

func TestStatementAveragesAndMissingGrades(t *testing.T) {
	data := statementData{
		Students: []StudentResult{{ID: 1}, {ID: 2}, {ID: 3}},
		Subjects: []subject.Subject{{ID: 1, Name: "Базы данных"}, {ID: 2, Name: "Физика"}},
		Grades:   []gradeRow{{1, 1, 5}, {1, 2, 4}, {2, 1, 3}},
	}
	result := buildStatement(data)
	if result.Average == nil || *result.Average != 4 || result.GradeCount != 3 {
		t.Fatalf("неверный средний балл группы: %+v", result)
	}
	if *result.Students[0].Average != 4.5 || *result.Students[1].Average != 3 {
		t.Fatal("неверный средний балл студента")
	}
	if result.Students[1].Grades[1].Value != nil || result.Students[2].Average != nil {
		t.Fatal("отсутствующая оценка не должна превращаться в ноль")
	}
	if *result.Subjects[0].Average != 4 || *result.Subjects[1].Average != 4 {
		t.Fatal("неверный средний балл по дисциплине")
	}
}

func TestStatementWithoutGrades(t *testing.T) {
	result := buildStatement(statementData{Students: []StudentResult{}, Subjects: []subject.Subject{}})
	if result.Average != nil || result.StudentCount != 0 || result.Subjects == nil || result.Students == nil {
		t.Fatalf("неверная пустая ведомость: %+v", result)
	}
}
