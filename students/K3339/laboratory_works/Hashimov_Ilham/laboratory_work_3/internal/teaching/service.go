package teaching

import (
	"context"
	"errors"
)

var ErrInvalidTeaching = errors.New("некорректные данные: teaching")

type TeachingService struct {
	repo *TeachingRepository
}

func NewTeachingService(repo *TeachingRepository) *TeachingService {
	return &TeachingService{repo: repo}
}

func validateTeaching(assignment *TeachingAssignment) error {
	if assignment == nil {
		return ErrInvalidTeaching
	}

	if assignment.TeacherID <= 0 ||
		assignment.SubjectID <= 0 ||
		assignment.GroupID <= 0 ||
		assignment.SemesterID <= 0 {
		return ErrInvalidTeaching
	}

	return nil
}

func (s *TeachingService) Create(ctx context.Context, assignment *TeachingAssignment) error {
	if err := validateTeaching(assignment); err != nil {
		return err
	}

	assignment.ID = 0
	assignment.Subject = nil

	return s.repo.Create(ctx, assignment)
}

func (s *TeachingService) GetByID(ctx context.Context, id int64) (TeachingAssignment, error) {
	if id <= 0 {
		return TeachingAssignment{}, ErrInvalidTeaching
	}

	return s.repo.GetByID(ctx, id)
}

func (s *TeachingService) GetAll(ctx context.Context) ([]TeachingAssignment, error) {
	return s.repo.GetAll(ctx)
}

func (s *TeachingService) Update(ctx context.Context, assignment TeachingAssignment) error {
	if assignment.ID <= 0 {
		return ErrInvalidTeaching
	}

	if err := validateTeaching(&assignment); err != nil {
		return err
	}

	assignment.Subject = nil
	return s.repo.Update(ctx, assignment)
}

func (s *TeachingService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidTeaching
	}

	return s.repo.Delete(ctx, id)
}
