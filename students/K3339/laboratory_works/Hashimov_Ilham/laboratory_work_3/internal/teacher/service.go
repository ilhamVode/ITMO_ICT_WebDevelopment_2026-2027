package teacher

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidTeacher = errors.New("некорректные данные: teacher")

type TeacherService struct {
	repo *TeacherRepository
}

func NewTeacherService(repo *TeacherRepository) *TeacherService {
	return &TeacherService{repo: repo}
}

func validateTeacher(teacher *Teacher) error {
	if teacher == nil {
		return ErrInvalidTeacher
	}

	teacher.FirstName = strings.TrimSpace(teacher.FirstName)
	teacher.LastName = strings.TrimSpace(teacher.LastName)

	if teacher.FirstName == "" || teacher.LastName == "" {
		return ErrInvalidTeacher
	}

	if teacher.FixedClassroomID != nil && *teacher.FixedClassroomID <= 0 {
		return ErrInvalidTeacher
	}

	if teacher.MiddleName != nil {
		value := strings.TrimSpace(*teacher.MiddleName)

		if value == "" {
			teacher.MiddleName = nil
		} else {
			teacher.MiddleName = &value
		}
	}

	return nil
}

func (s *TeacherService) Create(ctx context.Context, teacher *Teacher) error {
	if err := validateTeacher(teacher); err != nil {
		return err
	}

	teacher.ID = 0
	teacher.IsActive = true
	teacher.Assignments = nil

	return s.repo.Create(ctx, teacher)
}

func (s *TeacherService) GetByID(ctx context.Context, id int64) (Teacher, error) {
	if id <= 0 {
		return Teacher{}, ErrInvalidTeacher
	}

	return s.repo.GetByID(ctx, id)
}

func (s *TeacherService) GetAll(ctx context.Context) ([]Teacher, error) {
	return s.repo.GetAll(ctx)
}

func (s *TeacherService) GetWithSubjects(ctx context.Context, id int64) (Teacher, error) {
	if id <= 0 {
		return Teacher{}, ErrInvalidTeacher
	}

	return s.repo.GetWithSubjects(ctx, id)
}

func (s *TeacherService) Update(ctx context.Context, teacher Teacher) error {
	if teacher.ID <= 0 {
		return ErrInvalidTeacher
	}

	if err := validateTeacher(&teacher); err != nil {
		return err
	}

	teacher.Assignments = nil
	return s.repo.Update(ctx, teacher)
}

func (s *TeacherService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidTeacher
	}

	return s.repo.Delete(ctx, id)
}

// Увольнение преподавателя без удаления истории
func (s *TeacherService) Dismiss(ctx context.Context, id int64) error {
	teacher, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}

	teacher.IsActive = false
	return s.repo.Update(ctx, teacher)
}
