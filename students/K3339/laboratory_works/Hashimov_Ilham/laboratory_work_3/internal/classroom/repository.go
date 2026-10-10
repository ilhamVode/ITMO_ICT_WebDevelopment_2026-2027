package classroom

import (
	"context"

	"gorm.io/gorm"
)

type ClassroomRepository struct {
	db *gorm.DB
}

func NewClassroomRepository(db *gorm.DB) *ClassroomRepository {
	return &ClassroomRepository{db: db}
}

func (r *ClassroomRepository) Create(ctx context.Context, classroom *Classroom) error {
	return r.db.WithContext(ctx).Create(classroom).Error
}

func (r *ClassroomRepository) GetByID(ctx context.Context, id int64) (Classroom, error) {
	var classroom Classroom

	err := r.db.WithContext(ctx).First(&classroom, id).Error
	return classroom, err
}

func (r *ClassroomRepository) GetAll(ctx context.Context) ([]Classroom, error) {
	var classrooms []Classroom

	err := r.db.WithContext(ctx).Order("id ASC").Find(&classrooms).Error
	return classrooms, err
}

func (r *ClassroomRepository) Update(ctx context.Context, classroom Classroom) error {
	result := r.db.WithContext(ctx).
		Model(&Classroom{}).
		Where("id = ?", classroom.ID).
		Select("number", "building").
		Updates(&classroom)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *ClassroomRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Classroom{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
