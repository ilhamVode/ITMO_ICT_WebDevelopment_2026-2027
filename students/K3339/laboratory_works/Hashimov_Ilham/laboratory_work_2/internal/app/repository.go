package app

import (
	"context"
)

type UserRepository interface {
	Create(ctx context.Context, user User) (int64, error)

	GetByID(ctx context.Context, id int64) (User, error)
	GetByUsername(ctx context.Context, username string) (User, error)
	GetAll(ctx context.Context) ([]User, error)
}

type HomeworkRepository interface {
	Create(ctx context.Context, homework Homework) (int64, error)

	GetByID(ctx context.Context, id int64) (Homework, error)
	GetAll(ctx context.Context) ([]Homework, error)

	GetPage(ctx context.Context, limit int, offset int) ([]Homework, error)
	Count(ctx context.Context) (int64, error)

	GetPageFiltered(ctx context.Context, limit int, offset int, search string, subject string) ([]Homework, error)
	CountFiltered(ctx context.Context, search string, subject string) (int64, error)

	Update(ctx context.Context, homework Homework) error
	Delete(ctx context.Context, id int64) error
}

type SubmissionRepository interface {
	Create(ctx context.Context, submission Submission) (int64, error)

	GetByID(ctx context.Context, id int64) (Submission, error)
	GetByStudentID(ctx context.Context, studentID int64) ([]Submission, error)
	GetByHomeworkID(ctx context.Context, homeworkID int64) ([]Submission, error)
	GetAll(ctx context.Context) ([]Submission, error)

	Update(ctx context.Context, submission Submission) error
	Delete(ctx context.Context, id int64) error

	SetGrade(ctx context.Context, submissionID int64, grade float64) error
}
