package teaching

import (
	"context"

	"gorm.io/gorm"
)

type TeachingRepository struct {
	db *gorm.DB
}

func NewTeachingRepository(db *gorm.DB) *TeachingRepository {
	return &TeachingRepository{db: db}
}

func (r *TeachingRepository) Create(ctx context.Context, assignment *TeachingAssignment) error {
	return r.db.WithContext(ctx).Create(assignment).Error
}

func (r *TeachingRepository) GetByID(ctx context.Context, id int64) (TeachingAssignment, error) {
	var assignment TeachingAssignment

	err := r.db.WithContext(ctx).First(&assignment, id).Error
	return assignment, err
}

func (r *TeachingRepository) GetAll(ctx context.Context) ([]TeachingAssignment, error) {
	var assignments []TeachingAssignment

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&assignments).Error

	return assignments, err
}

func (r *TeachingRepository) Update(ctx context.Context, assignment TeachingAssignment) error {
	result := r.db.WithContext(ctx).
		Model(&TeachingAssignment{}).
		Where("id = ?", assignment.ID).
		Select(
			"teacher_id",
			"subject_id",
			"group_id",
			"semester_id",
		).
		Updates(&assignment)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TeachingRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&TeachingAssignment{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
