package group

import (
	"context"
	"errors"
	"strings"
)

var ErrInvalidGroup = errors.New("некорректные данные: group")

type GroupService struct {
	repo *GroupRepository
}

func NewGroupService(repo *GroupRepository) *GroupService {
	return &GroupService{repo: repo}
}

func validateGroup(group *StudyGroup) error {
	if group == nil {
		return ErrInvalidGroup
	}

	group.Name = strings.TrimSpace(group.Name)

	if group.Name == "" || group.Course <= 0 {
		return ErrInvalidGroup
	}

	return nil
}

func (s *GroupService) Create(ctx context.Context, group *StudyGroup) error {
	if err := validateGroup(group); err != nil {
		return err
	}

	group.ID = 0
	group.Students = nil

	return s.repo.Create(ctx, group)
}

func (s *GroupService) GetByID(ctx context.Context, id int64) (StudyGroup, error) {
	if id <= 0 {
		return StudyGroup{}, ErrInvalidGroup
	}

	return s.repo.GetByID(ctx, id)
}

func (s *GroupService) GetAll(ctx context.Context) ([]StudyGroup, error) {
	return s.repo.GetAll(ctx)
}

func (s *GroupService) GetWithStudents(ctx context.Context, id int64) (StudyGroup, error) {
	if id <= 0 {
		return StudyGroup{}, ErrInvalidGroup
	}

	return s.repo.GetWithStudents(ctx, id)
}

func (s *GroupService) Update(ctx context.Context, group StudyGroup) error {
	if group.ID <= 0 {
		return ErrInvalidGroup
	}

	if err := validateGroup(&group); err != nil {
		return err
	}

	group.Students = nil
	return s.repo.Update(ctx, group)
}

func (s *GroupService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidGroup
	}

	return s.repo.Delete(ctx, id)
}
