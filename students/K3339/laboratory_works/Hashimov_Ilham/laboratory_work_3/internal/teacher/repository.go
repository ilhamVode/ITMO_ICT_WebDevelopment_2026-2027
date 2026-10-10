package teacher

import (
	"context"

	"gorm.io/gorm"
)

type TeacherRepository struct {
	db *gorm.DB
}

func NewTeacherRepository(db *gorm.DB) *TeacherRepository {
	return &TeacherRepository{db: db}
}

func (r *TeacherRepository) Create(ctx context.Context, teacher *Teacher) error {
	return r.db.WithContext(ctx).Create(teacher).Error
}

func (r *TeacherRepository) GetByID(ctx context.Context, id int64) (Teacher, error) {
	var teacher Teacher

	err := r.db.WithContext(ctx).First(&teacher, id).Error
	return teacher, err
}

func (r *TeacherRepository) GetAll(ctx context.Context) ([]Teacher, error) {
	var teachers []Teacher

	err := r.db.WithContext(ctx).Order("id ASC").Find(&teachers).Error
	return teachers, err
}

// Получение преподавателя с назначениями и дисциплинами
func (r *TeacherRepository) GetWithSubjects(ctx context.Context, id int64) (Teacher, error) {
	var teacher Teacher

	err := r.db.WithContext(ctx).
		Preload("Assignments.Subject").
		First(&teacher, id).Error

	// Один преподаватель может вести дисциплину в нескольких группах.
	// Во вложенном списке дисциплин убираем такие повторы.
	seen := make(map[int64]bool)
	for _, assignment := range teacher.Assignments {
		if assignment.Subject != nil && !seen[assignment.SubjectID] {
			teacher.Subjects = append(teacher.Subjects, *assignment.Subject)
			seen[assignment.SubjectID] = true
		}
	}
	return teacher, err
}

func (r *TeacherRepository) Update(ctx context.Context, teacher Teacher) error {
	result := r.db.WithContext(ctx).
		Model(&Teacher{}).
		Where("id = ?", teacher.ID).
		Select(
			"first_name",
			"last_name",
			"middle_name",
			"fixed_classroom_id",
			"is_active",
		).
		Updates(&teacher)

	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

func (r *TeacherRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM schedules WHERE assignment_id IN (SELECT id FROM teaching_assignments WHERE teacher_id = ?)", id).Error; err != nil {
			return err
		}
		if err := tx.Exec("DELETE FROM teaching_assignments WHERE teacher_id = ?", id).Error; err != nil {
			return err
		}
		result := tx.Delete(&Teacher{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
