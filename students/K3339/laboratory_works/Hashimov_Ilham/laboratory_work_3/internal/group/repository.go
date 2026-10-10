package group

import (
	"context"

	"gorm.io/gorm"
)

type GroupRepository struct {
	db *gorm.DB
}

func NewGroupRepository(db *gorm.DB) *GroupRepository {
	return &GroupRepository{db: db}
}

func (r *GroupRepository) Create(ctx context.Context, group *StudyGroup) error {
	return r.db.WithContext(ctx).Create(group).Error
}

func (r *GroupRepository) GetByID(ctx context.Context, id int64) (StudyGroup, error) {
	var group StudyGroup

	err := r.db.WithContext(ctx).First(&group, id).Error
	return group, err
}

func (r *GroupRepository) GetAll(ctx context.Context) ([]StudyGroup, error) {
	var groups []StudyGroup

	err := r.db.WithContext(ctx).Order("id ASC").Find(&groups).Error
	return groups, err
}

// Получение группы вместе со студентами
func (r *GroupRepository) GetWithStudents(ctx context.Context, id int64) (StudyGroup, error) {
	var group StudyGroup

	err := r.db.WithContext(ctx).
		Preload("Students").
		First(&group, id).Error

	return group, err
}

func (r *GroupRepository) Update(ctx context.Context, group StudyGroup) error {
	result := r.db.WithContext(ctx).
		Model(&StudyGroup{}).
		Where("id = ?", group.ID).
		Select("name", "course").
		Updates(&group)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *GroupRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&StudyGroup{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
