package subject

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidSubject = errors.New("некорректные данные: subject")

type SubjectService struct {
	repo *SubjectRepository
}

func NewSubjectService(repo *SubjectRepository) *SubjectService {
	return &SubjectService{repo: repo}
}

func validateSubject(subject *Subject) error {
	if subject == nil {
		return ErrInvalidSubject
	}

	subject.Name = strings.TrimSpace(subject.Name)

	if subject.Name == "" {
		return ErrInvalidSubject
	}

	return nil
}

func (s *SubjectService) Create(ctx context.Context, subject *Subject) error {
	if err := validateSubject(subject); err != nil {
		return err
	}

	subject.ID = 0
	return s.repo.Create(ctx, subject)
}

func (s *SubjectService) GetByID(ctx context.Context, id int64) (Subject, error) {
	if id <= 0 {
		return Subject{}, ErrInvalidSubject
	}

	return s.repo.GetByID(ctx, id)
}

func (s *SubjectService) GetAll(ctx context.Context) ([]Subject, error) {
	return s.repo.GetAll(ctx)
}

func (s *SubjectService) Update(ctx context.Context, subject Subject) error {
	if subject.ID <= 0 {
		return ErrInvalidSubject
	}

	if err := validateSubject(&subject); err != nil {
		return err
	}

	return s.repo.Update(ctx, subject)
}

func (s *SubjectService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidSubject
	}

	return s.repo.Delete(ctx, id)
}
