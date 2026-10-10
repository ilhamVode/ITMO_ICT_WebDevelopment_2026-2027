package report

import (
	"context"
	"errors"
	"math"
)

var ErrInvalidReport = errors.New("необходимо указать положительные group_id и semester_id")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func average(sum, count int) *float64 {
	if count == 0 {
		return nil
	}
	value := math.Round(float64(sum)/float64(count)*100) / 100
	return &value
}

func buildStatement(data statementData) Statement {
	result := Statement{
		Group: data.Group, Semester: data.Semester,
		Students: data.Students, Subjects: make([]SubjectResult, 0),
		StudentCount: len(data.Students),
	}
	// Ключ состоит из ID студента и ID дисциплины.
	values := make(map[[2]int64]int)
	subjectSums := make(map[int64]int)
	subjectCounts := make(map[int64]int)
	for _, grade := range data.Grades {
		values[[2]int64{grade.StudentID, grade.SubjectID}] = grade.Value
	}
	total := 0
	for i := range result.Students {
		student := &result.Students[i]
		student.Grades = make([]StudentGrade, 0, len(data.Subjects))
		sum, count := 0, 0
		for _, subject := range data.Subjects {
			grade := StudentGrade{SubjectID: subject.ID, SubjectName: subject.Name}
			if value, ok := values[[2]int64{student.ID, subject.ID}]; ok {
				grade.Value = &value
				sum += value
				count++
				total += value
				result.GradeCount++
				subjectSums[subject.ID] += value
				subjectCounts[subject.ID]++
			}
			student.Grades = append(student.Grades, grade)
		}
		student.Average = average(sum, count)
	}
	for _, subject := range data.Subjects {
		result.Subjects = append(result.Subjects, SubjectResult{
			ID: subject.ID, Name: subject.Name, GradeCount: subjectCounts[subject.ID],
			Average: average(subjectSums[subject.ID], subjectCounts[subject.ID]),
		})
	}
	// Средний балл группы считаем по всем выставленным оценкам.
	result.Average = average(total, result.GradeCount)
	return result
}

func (s *Service) GetStatement(ctx context.Context, groupID, semesterID int64) (Statement, error) {
	if groupID <= 0 || semesterID <= 0 {
		return Statement{}, ErrInvalidReport
	}
	data, err := s.repo.GetStatement(ctx, groupID, semesterID)
	if err != nil {
		return Statement{}, err
	}
	return buildStatement(data), nil
}
