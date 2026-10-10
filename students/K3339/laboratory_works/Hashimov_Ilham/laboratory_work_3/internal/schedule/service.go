package schedule

import (
	"context"
	"errors"
)

var ErrInvalidSchedule = errors.New("некорректные данные: schedule")

type ScheduleService struct {
	repo *ScheduleRepository
}

func NewScheduleService(repo *ScheduleRepository) *ScheduleService {
	return &ScheduleService{repo: repo}
}

func validateSchedule(schedule *Schedule) error {
	if schedule == nil {
		return ErrInvalidSchedule
	}

	if schedule.AssignmentID <= 0 ||
		schedule.ClassroomID <= 0 ||
		schedule.Weekday < 1 ||
		schedule.Weekday > 7 ||
		schedule.LessonNumber <= 0 {
		return ErrInvalidSchedule
	}

	return nil
}

func (s *ScheduleService) Create(ctx context.Context, schedule *Schedule) error {
	if err := validateSchedule(schedule); err != nil {
		return err
	}

	schedule.ID = 0
	schedule.Assignment = nil
	schedule.Classroom = nil

	return s.repo.Create(ctx, schedule)
}

func (s *ScheduleService) GetByID(ctx context.Context, id int64) (Schedule, error) {
	if id <= 0 {
		return Schedule{}, ErrInvalidSchedule
	}

	return s.repo.GetByID(ctx, id)
}

func (s *ScheduleService) GetAll(ctx context.Context) ([]Schedule, error) {
	return s.repo.GetAll(ctx)
}

func (s *ScheduleService) GetWithDetails(ctx context.Context, id int64) (Schedule, error) {
	if id <= 0 {
		return Schedule{}, ErrInvalidSchedule
	}

	return s.repo.GetWithDetails(ctx, id)
}

func (s *ScheduleService) Update(ctx context.Context, schedule Schedule) error {
	if schedule.ID <= 0 {
		return ErrInvalidSchedule
	}

	if err := validateSchedule(&schedule); err != nil {
		return err
	}

	schedule.Assignment = nil
	schedule.Classroom = nil

	return s.repo.Update(ctx, schedule)
}

func (s *ScheduleService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidSchedule
	}

	return s.repo.Delete(ctx, id)
}
