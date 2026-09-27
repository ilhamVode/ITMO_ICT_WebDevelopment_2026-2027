package app

import (
	"context"

	"gorm.io/gorm"
)

type PostgresUserRepository struct {
	db *gorm.DB
}

type PostgresHomeworkRepository struct {
	db *gorm.DB
}

type PostgresSubmissionRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func NewPostgresHomeworkRepository(db *gorm.DB) *PostgresHomeworkRepository {
	return &PostgresHomeworkRepository{db: db}
}

func NewPostgresSubmissionRepository(db *gorm.DB) *PostgresSubmissionRepository {
	return &PostgresSubmissionRepository{db: db}
}

func (r *PostgresUserRepository) Create(ctx context.Context, user User) (int64, error) {

	result := r.db.WithContext(ctx).Create(&user)

	if result.Error != nil {
		return 0, result.Error
	}

	return user.ID, nil
}

func (r *PostgresUserRepository) GetByID(ctx context.Context, id int64) (User, error) {
	user := User{}

	result := r.db.WithContext(ctx).First(&user, id)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (r *PostgresUserRepository) GetByUsername(ctx context.Context, username string) (User, error) {
	user := User{}

	result := r.db.WithContext(ctx).Where(&User{Username: username}).First(&user)

	if result.Error != nil {
		return user, result.Error
	}

	return user, nil
}

func (r *PostgresUserRepository) GetAll(ctx context.Context) ([]User, error) {
	var users []User

	result := r.db.WithContext(ctx).Find(&users)

	if result.Error != nil {
		return users, result.Error
	}

	return users, nil
}

func (r *PostgresHomeworkRepository) Create(ctx context.Context, homework Homework) (int64, error) {
	result := r.db.WithContext(ctx).Create(&homework)

	if result.Error != nil {
		return 0, result.Error
	}

	return homework.ID, nil
}

func (r *PostgresHomeworkRepository) GetByID(ctx context.Context, id int64) (Homework, error) {
	homework := Homework{}

	result := r.db.WithContext(ctx).First(&homework, id)

	if result.Error != nil {
		return homework, result.Error
	}

	return homework, nil
}

func (r *PostgresHomeworkRepository) GetAll(ctx context.Context) ([]Homework, error) {
	var homeworks []Homework

	result := r.db.WithContext(ctx).Find(&homeworks)

	if result.Error != nil {
		return homeworks, result.Error
	}

	return homeworks, nil
}

func (r *PostgresHomeworkRepository) GetPage(ctx context.Context, limit int, offset int) ([]Homework, error) {
	var homeworks []Homework

	err := r.db.WithContext(ctx).
		Limit(limit).
		Offset(offset).
		Find(&homeworks).Error

	if err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (r *PostgresHomeworkRepository) Count(ctx context.Context) (int64, error) {
	var count int64

	err := r.db.WithContext(ctx).
		Model(&Homework{}).
		Count(&count).Error

	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresHomeworkRepository) GetPageFiltered(ctx context.Context, limit int, offset int, search string, subject string) ([]Homework, error) {
	var homeworks []Homework

	query := r.db.WithContext(ctx).Model(&Homework{})

	if search != "" {
		searchValue := "%" + search + "%"

		query = query.Where(
			"subject ILIKE ? OR text ILIKE ?",
			searchValue,
			searchValue,
		)
	}

	if subject != "" {
		query = query.Where("subject = ?", subject)
	}

	err := query.
		Limit(limit).
		Offset(offset).
		Find(&homeworks).Error

	if err != nil {
		return nil, err
	}

	return homeworks, nil
}

func (r *PostgresHomeworkRepository) CountFiltered(ctx context.Context, search string, subject string) (int64, error) {
	var count int64

	query := r.db.WithContext(ctx).Model(&Homework{})

	if search != "" {
		searchValue := "%" + search + "%"

		query = query.Where(
			"subject ILIKE ? OR text ILIKE ?",
			searchValue,
			searchValue,
		)
	}

	if subject != "" {
		query = query.Where("subject = ?", subject)
	}

	err := query.Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (r *PostgresHomeworkRepository) Update(ctx context.Context, homework Homework) error {
	_, err := r.GetByID(ctx, homework.ID)

	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Save(&homework)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresHomeworkRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.GetByID(ctx, id)

	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Delete(&Homework{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresSubmissionRepository) Create(ctx context.Context, submission Submission) (int64, error) {
	result := r.db.WithContext(ctx).Create(&submission)

	if result.Error != nil {
		return 0, result.Error
	}

	return submission.ID, nil
}

func (r *PostgresSubmissionRepository) GetByID(ctx context.Context, id int64) (Submission, error) {
	submission := Submission{}

	result := r.db.WithContext(ctx).First(&submission, id)

	if result.Error != nil {
		return submission, result.Error
	}

	return submission, nil
}

func (r *PostgresSubmissionRepository) GetByStudentID(ctx context.Context, studentID int64) ([]Submission, error) {
	submission := []Submission{}

	result := r.db.WithContext(ctx).Where(&Submission{StudentID: studentID}).Find(&submission)

	if result.Error != nil {
		return submission, result.Error
	}

	return submission, nil
}

func (r *PostgresSubmissionRepository) GetByHomeworkID(ctx context.Context, homeworkID int64) ([]Submission, error) {
	submission := []Submission{}

	result := r.db.WithContext(ctx).Where(&Submission{HomeworkID: homeworkID}).Find(&submission)

	if result.Error != nil {
		return submission, result.Error
	}

	return submission, nil
}

func (r *PostgresSubmissionRepository) Update(ctx context.Context, submission Submission) error {
	_, err := r.GetByID(ctx, submission.ID)

	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Save(&submission)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresSubmissionRepository) Delete(ctx context.Context, id int64) error {
	_, err := r.GetByID(ctx, id)

	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Delete(&Submission{}, id)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresSubmissionRepository) SetGrade(ctx context.Context, submissionID int64, grade float64) error {
	_, err := r.GetByID(ctx, submissionID)

	if err != nil {
		return err
	}

	result := r.db.WithContext(ctx).Model(&Submission{}).Where(&Submission{ID: submissionID}).Update("grade", grade)

	if result.Error != nil {
		return result.Error
	}

	return nil
}

func (r *PostgresSubmissionRepository) GetAll(ctx context.Context) ([]Submission, error) {
	var submissions []Submission

	result := r.db.WithContext(ctx).Find(&submissions)

	if result.Error != nil {
		return submissions, result.Error
	}

	return submissions, nil
}
