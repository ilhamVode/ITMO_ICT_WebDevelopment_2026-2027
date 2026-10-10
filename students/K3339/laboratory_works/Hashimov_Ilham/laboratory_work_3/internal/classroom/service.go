package classroom

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidClassroom = errors.New("некорректные данные: classroom")

type ClassroomService struct {
	repo *ClassroomRepository
}

func NewClassroomService(repo *ClassroomRepository) *ClassroomService {
	return &ClassroomService{repo: repo}
}

func validateClassroom(classroom *Classroom) error {
	if classroom == nil {
		return ErrInvalidClassroom
	}

	classroom.Number = strings.TrimSpace(classroom.Number)
	classroom.Building = strings.TrimSpace(classroom.Building)

	if classroom.Number == "" || classroom.Building == "" {
		return ErrInvalidClassroom
	}

	return nil
}

func (s *ClassroomService) Create(ctx context.Context, classroom *Classroom) error {
	if err := validateClassroom(classroom); err != nil {
		return err
	}

	classroom.ID = 0
	return s.repo.Create(ctx, classroom)
}

func (s *ClassroomService) GetByID(ctx context.Context, id int64) (Classroom, error) {
	if id <= 0 {
		return Classroom{}, ErrInvalidClassroom
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ClassroomService) GetAll(ctx context.Context) ([]Classroom, error) {
	return s.repo.GetAll(ctx)
}

func (s *ClassroomService) Update(ctx context.Context, classroom Classroom) error {
	if classroom.ID <= 0 {
		return ErrInvalidClassroom
	}

	if err := validateClassroom(&classroom); err != nil {
		return err
	}

	return s.repo.Update(ctx, classroom)
}

func (s *ClassroomService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidClassroom
	}

	return s.repo.Delete(ctx, id)
}
