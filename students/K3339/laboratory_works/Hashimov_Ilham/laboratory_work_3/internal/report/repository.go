package report

import (
	"context"
	"database/sql"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) GetStatement(ctx context.Context, groupID, semesterID int64) (statementData, error) {
	data := statementData{
		Students: make([]StudentResult, 0),
		Grades:   make([]gradeRow, 0),
	}
	// Все части ведомости читаем из одного снимка базы.
	err := r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		if err := tx.First(&data.Group, groupID).Error; err != nil {
			return err
		}
		if err := tx.First(&data.Semester, semesterID).Error; err != nil {
			return err
		}
		if err := tx.Table("students").Select("id, first_name, last_name, is_active").
			Where("group_id = ?", groupID).
			Order("last_name, first_name, id").Scan(&data.Students).Error; err != nil {
			return err
		}
		// Показываем и назначенные дисциплины без оценок, и сохранённые оценки.
		if err := tx.Raw(`
			SELECT id, name FROM subjects WHERE id IN (
				SELECT subject_id FROM teaching_assignments WHERE group_id = ? AND semester_id = ?
				UNION
				SELECT g.subject_id FROM grades g JOIN students st ON st.id = g.student_id
				WHERE st.group_id = ? AND g.semester_id = ?
			) ORDER BY name, id
		`, groupID, semesterID, groupID, semesterID).Scan(&data.Subjects).Error; err != nil {
			return err
		}
		return tx.Raw(`
			SELECT g.student_id, g.subject_id, g.value
			FROM grades g JOIN students st ON st.id = g.student_id
			WHERE st.group_id = ? AND g.semester_id = ?
		`, groupID, semesterID).Scan(&data.Grades).Error
	}, &sql.TxOptions{Isolation: sql.LevelRepeatableRead, ReadOnly: true})
	return data, err
}
