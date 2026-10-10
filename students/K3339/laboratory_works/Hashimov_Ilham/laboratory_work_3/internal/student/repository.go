package student

import (
	"context"

	"gorm.io/gorm"
)

type StudentRepository struct {
	db *gorm.DB
}

func NewStudentRepository(db *gorm.DB) *StudentRepository {
	return &StudentRepository{db: db}
}

// Создание студента
func (r *StudentRepository) Create(ctx context.Context, student *Student) error {
	return r.db.WithContext(ctx).Create(student).Error
}

// Получение студента по ID
func (r *StudentRepository) GetByID(ctx context.Context, id int64) (Student, error) {
	var student Student

	err := r.db.WithContext(ctx).First(&student, id).Error
	return student, err
}

// Получение всех студентов
func (r *StudentRepository) GetAll(ctx context.Context) ([]Student, error) {
	var students []Student

	err := r.db.WithContext(ctx).
		Order("id ASC").
		Find(&students).Error

	return students, err
}

// Обновление студента
func (r *StudentRepository) Update(ctx context.Context, student Student) error {
	result := r.db.WithContext(ctx).
		Model(&Student{}).
		Where("id = ?", student.ID).
		Updates(map[string]interface{}{
			"first_name": student.FirstName,
			"last_name":  student.LastName,
			"birth_date": student.BirthDate,
			"group_id":   student.GroupID,
			"is_active":  student.IsActive,
		})

	if result.Error != nil {
		return result.Error
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}

	return nil
}

// Удаление студента
func (r *StudentRepository) Delete(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.Exec("DELETE FROM grades WHERE student_id = ?", id).Error; err != nil {
			return err
		}
		result := tx.Delete(&Student{}, id)
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return gorm.ErrRecordNotFound
		}
		return nil
	})
}
