package schedule

import (
	"context"

	"gorm.io/gorm"
)

type ScheduleRepository struct {
	db *gorm.DB
}

func NewScheduleRepository(db *gorm.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

func (r *ScheduleRepository) Create(ctx context.Context, schedule *Schedule) error {
	return r.db.WithContext(ctx).Create(schedule).Error
}

func (r *ScheduleRepository) GetByID(ctx context.Context, id int64) (Schedule, error) {
	var schedule Schedule

	err := r.db.WithContext(ctx).First(&schedule, id).Error
	return schedule, err
}

func (r *ScheduleRepository) GetAll(ctx context.Context) ([]Schedule, error) {
	var schedules []Schedule

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&schedules).Error

	return schedules, err
}

// Получаем занятие вместе с назначением, дисциплиной и кабинетом
func (r *ScheduleRepository) GetWithDetails(ctx context.Context, id int64) (Schedule, error) {
	var schedule Schedule

	err := r.db.WithContext(ctx).
		Preload("Assignment.Subject").
		Preload("Classroom").
		First(&schedule, id).Error

	return schedule, err
}

func (r *ScheduleRepository) Update(ctx context.Context, schedule Schedule) error {
	result := r.db.WithContext(ctx).
		Model(&Schedule{}).
		Where("id = ?", schedule.ID).
		Select(
			"assignment_id",
			"classroom_id",
			"weekday",
			"lesson_number",
		).
		Updates(&schedule)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ScheduleRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Schedule{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
