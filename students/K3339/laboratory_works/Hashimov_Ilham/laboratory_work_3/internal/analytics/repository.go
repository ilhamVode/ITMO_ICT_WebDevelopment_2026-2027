package analytics

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// 1 и 4. Получение занятий группы на выбранный день.
// Если lessonNumber == nil, получаем расписание на весь день.
func (r *Repository) getLessons(
	ctx context.Context,
	groupID int64,
	semesterID int64,
	weekday int,
	lessonNumber *int,
) ([]Lesson, error) {
	lessons := make([]Lesson, 0)

	query := `
		SELECT
			sc.id AS schedule_id,
			sc.lesson_number,
			sub.id AS subject_id,
			sub.name AS subject_name,
			t.id AS teacher_id,
			t.first_name AS teacher_first_name,
			t.last_name AS teacher_last_name,
			c.id AS classroom_id,
			c.number AS classroom_number,
			c.building AS classroom_building
		FROM schedules sc
		JOIN teaching_assignments ta ON ta.id = sc.assignment_id
		JOIN subjects sub ON sub.id = ta.subject_id
		JOIN teachers t ON t.id = ta.teacher_id
		JOIN classrooms c ON c.id = sc.classroom_id
		WHERE ta.group_id = ?
		  AND ta.semester_id = ?
		  AND sc.weekday = ?
	`

	args := []any{groupID, semesterID, weekday}

	if lessonNumber != nil {
		query += " AND sc.lesson_number = ?"
		args = append(args, *lessonNumber)
	}

	query += " ORDER BY sc.lesson_number, sc.id"

	err := r.db.WithContext(ctx).
		Raw(query, args...).
		Scan(&lessons).Error

	return lessons, err
}

// 1. Какой предмет у группы на конкретной паре?
func (r *Repository) GetLesson(
	ctx context.Context,
	groupID, semesterID int64,
	weekday, lessonNumber int,
) ([]Lesson, error) {
	return r.getLessons(ctx, groupID, semesterID, weekday, &lessonNumber)
}

// 2. Какие преподаватели ведут занятия у группы?
func (r *Repository) GetTeachers(
	ctx context.Context,
	groupID, semesterID int64,
) ([]TeacherResult, error) {
	teachers := make([]TeacherResult, 0)

	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT
			t.id,
			t.first_name,
			t.last_name
		FROM teachers t
		JOIN teaching_assignments ta ON ta.teacher_id = t.id
		WHERE ta.group_id = ?
		  AND ta.semester_id = ?
		ORDER BY t.last_name, t.first_name, t.id
	`, groupID, semesterID).Scan(&teachers).Error

	return teachers, err
}

// 3. Каким группам преподаватель ведёт выбранный предмет?
func (r *Repository) GetGroups(
	ctx context.Context,
	teacherID, subjectID, semesterID int64,
) ([]GroupResult, error) {
	groups := make([]GroupResult, 0)

	err := r.db.WithContext(ctx).Raw(`
		SELECT DISTINCT
			g.id,
			g.name,
			g.course
		FROM study_groups g
		JOIN teaching_assignments ta ON ta.group_id = g.id
		WHERE ta.teacher_id = ?
		  AND ta.subject_id = ?
		  AND ta.semester_id = ?
		ORDER BY g.id
	`, teacherID, subjectID, semesterID).Scan(&groups).Error

	return groups, err
}

// 4. Полное расписание группы на выбранный день
func (r *Repository) GetSchedule(
	ctx context.Context,
	groupID, semesterID int64,
	weekday int,
) ([]Lesson, error) {
	return r.getLessons(ctx, groupID, semesterID, weekday, nil)
}

// 5. Количество активных студентов по курсам,
// у которых есть занятия в указанном кабинете за семестр
func (r *Repository) GetStudentCounts(
	ctx context.Context,
	classroomID, semesterID int64,
) ([]CourseCount, error) {
	counts := make([]CourseCount, 0)

	err := r.db.WithContext(ctx).Raw(`
		SELECT
			g.course,
			COUNT(DISTINCT st.id) AS student_count
		FROM schedules sc
		JOIN teaching_assignments ta ON ta.id = sc.assignment_id
		JOIN study_groups g ON g.id = ta.group_id
		JOIN students st ON st.group_id = g.id
		WHERE sc.classroom_id = ?
		  AND ta.semester_id = ?
		  AND st.is_active = TRUE
		GROUP BY g.course
		ORDER BY g.course
	`, classroomID, semesterID).Scan(&counts).Error

	return counts, err
}
