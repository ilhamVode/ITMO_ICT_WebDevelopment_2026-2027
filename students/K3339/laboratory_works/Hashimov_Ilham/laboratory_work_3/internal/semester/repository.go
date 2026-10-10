package semester

import (
	"context"

	"gorm.io/gorm"
)

type SemesterRepository struct {
	db *gorm.DB
}

func NewSemesterRepository(db *gorm.DB) *SemesterRepository {
	return &SemesterRepository{db: db}
}

func (r *SemesterRepository) Create(ctx context.Context, semester *Semester) error {
	return r.db.WithContext(ctx).Create(semester).Error
}

func (r *SemesterRepository) GetByID(ctx context.Context, id int64) (Semester, error) {
	var semester Semester

	err := r.db.WithContext(ctx).First(&semester, id).Error
	return semester, err
}

func (r *SemesterRepository) GetAll(ctx context.Context) ([]Semester, error) {
	var semesters []Semester

	err := r.db.WithContext(ctx).Order("id ASC").Find(&semesters).Error
	return semesters, err
}

func (r *SemesterRepository) Update(ctx context.Context, semester Semester) error {
	result := r.db.WithContext(ctx).
		Model(&Semester{}).
		Where("id = ?", semester.ID).
		Select("academic_year", "number").
		Updates(&semester)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *SemesterRepository) Delete(ctx context.Context, id int64) error {
	result := r.db.WithContext(ctx).Delete(&Semester{}, id)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}
