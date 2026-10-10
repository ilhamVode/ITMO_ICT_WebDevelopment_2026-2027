package subject

import (
	"context"

	"gorm.io/gorm"
)

type SubjectRepository struct {
	db *gorm.DB
}

func NewSubjectRepository(db *gorm.DB) *SubjectRepository {
	return &SubjectRepository{db: db}
}

func (r *SubjectRepository) Create(ctx context.Context, subject *Subject) error {
	return r.db.WithContext(ctx).Create(subject).Error
}

func (r *SubjectRepository) GetByID(ctx context.Context, id int64) (Subject, error) {
	var subject Subject

	err := r.db.WithContext(ctx).First(&subject, id).Error
	return subject, err
}

func (r *SubjectRepository) GetAll(ctx context.Context) ([]Subject, error) {
	var subjects []Subject

	err := r.db.WithContext(ctx).Order("id ASC").Find(&subjects).Error
	return subjects, err
}

func (r *SubjectRepository) Update(ctx context.Context, subject Subject) error {
	result := r.db.WithContext(ctx).
		Model(&Subject{}).
		Where("id = ?", subject.ID).
		Select("name").
		Updates(&subject)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SubjectRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Subject{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
