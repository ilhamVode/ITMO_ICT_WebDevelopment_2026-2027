package grade

import (
	"context"
	"errors"
)

var ErrInvalidGrade = errors.New("некорректные данные: grade")

type GradeService struct {
	repo *GradeRepository
}

func NewGradeService(repo *GradeRepository) *GradeService {
	return &GradeService{repo: repo}
}

func validateGrade(grade *Grade) error {
	if grade == nil {
		return ErrInvalidGrade
	}

	if grade.StudentID <= 0 ||
		grade.SubjectID <= 0 ||
		grade.SemesterID <= 0 ||
		grade.Value < 2 ||
		grade.Value > 5 {
		return ErrInvalidGrade
	}

	return nil
}

func (s *GradeService) Create(ctx context.Context, grade *Grade) error {
	if err := validateGrade(grade); err != nil {
		return err
	}

	grade.ID = 0
	grade.Student = nil
	grade.Subject = nil
	grade.Semester = nil

	return s.repo.Create(ctx, grade)
}

func (s *GradeService) GetByID(ctx context.Context, id int64) (Grade, error) {
	if id <= 0 {
		return Grade{}, ErrInvalidGrade
	}

	return s.repo.GetByID(ctx, id)
}

func (s *GradeService) GetAll(ctx context.Context) ([]Grade, error) {
	return s.repo.GetAll(ctx)
}

func (s *GradeService) GetWithDetails(ctx context.Context, id int64) (Grade, error) {
	if id <= 0 {
		return Grade{}, ErrInvalidGrade
	}

	return s.repo.GetWithDetails(ctx, id)
}

func (s *GradeService) Update(ctx context.Context, grade Grade) error {
	if grade.ID <= 0 {
		return ErrInvalidGrade
	}

	if err := validateGrade(&grade); err != nil {
		return err
	}

	grade.Student = nil
	grade.Subject = nil
	grade.Semester = nil

	return s.repo.Update(ctx, grade)
}

func (s *GradeService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidGrade
	}

	return s.repo.Delete(ctx, id)
}
