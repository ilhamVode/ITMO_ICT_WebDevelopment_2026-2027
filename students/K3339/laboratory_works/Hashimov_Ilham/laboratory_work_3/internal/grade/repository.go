package grade

import (
	"context"

	"gorm.io/gorm"
)

type GradeRepository struct {
	db *gorm.DB
}

func NewGradeRepository(db *gorm.DB) *GradeRepository {
	return &GradeRepository{db: db}
}

func (r *GradeRepository) Create(ctx context.Context, grade *Grade) error {
	return r.db.WithContext(ctx).Create(grade).Error
}

func (r *GradeRepository) GetByID(ctx context.Context, id int64) (Grade, error) {
	var grade Grade

	err := r.db.WithContext(ctx).First(&grade, id).Error
	return grade, err
}

func (r *GradeRepository) GetAll(ctx context.Context) ([]Grade, error) {
	var grades []Grade

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&grades).Error

	return grades, err
}

// Получаем оценку со студентом, дисциплиной и семестром
func (r *GradeRepository) GetWithDetails(ctx context.Context, id int64) (Grade, error) {
	var grade Grade

	err := r.db.WithContext(ctx).
		Preload("Student").
		Preload("Subject").
		Preload("Semester").
		First(&grade, id).Error

	return grade, err
}

func (r *GradeRepository) Update(ctx context.Context, grade Grade) error {
	result := r.db.WithContext(ctx).
		Model(&Grade{}).
		Where("id = ?", grade.ID).
		Select(
			"student_id",
			"subject_id",
			"semester_id",
			"value",
		).
		Updates(&grade)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GradeRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Grade{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
