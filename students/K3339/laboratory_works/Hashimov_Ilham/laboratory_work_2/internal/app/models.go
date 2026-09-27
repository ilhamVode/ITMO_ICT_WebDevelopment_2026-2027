package app

import (
	"time"
)

type User struct {
	ID           int64
	Username     string
	PasswordHash string
	IsAdmin      bool
	ClassName    string

	PassportNumber string
	HomeAddress    string
	Nationality    string
}

type Homework struct {
	ID        int64
	TeacherID int64
	Subject   string
	Text      string
	Penalty   string
	IssueDate time.Time
	StartDate time.Time
	EndDate   time.Time
}

type Submission struct {
	ID          int64
	StudentID   int64
	HomeworkID  int64
	Text        string
	SubmittedAt time.Time
	Grade       *float64
}
