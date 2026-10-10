package analytics

import (
	"context"
	"errors"
	"math"
)

var ErrInvalidQuery = errors.New("некорректные параметры запроса")

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func validIDs(ids ...int64) bool {
	for _, id := range ids {
		if id <= 0 {
			return false
		}
	}
	return true
}

func validWeekday(weekday int) bool {
	return weekday >= 1 && weekday <= 7
}

// Предмет на конкретной паре
func (s *Service) GetLesson(
	ctx context.Context,
	groupID, semesterID int64,
	weekday, lessonNumber int,
) ([]Lesson, error) {
	if !validIDs(groupID, semesterID) ||
		!validWeekday(weekday) ||
		lessonNumber <= 0 ||
		lessonNumber > math.MaxInt32 {
		return nil, ErrInvalidQuery
	}

	return s.repo.GetLesson(ctx, groupID, semesterID, weekday, lessonNumber)
}

// Преподаватели группы
func (s *Service) GetTeachers(
	ctx context.Context,
	groupID, semesterID int64,
) ([]TeacherResult, error) {
	if !validIDs(groupID, semesterID) {
		return nil, ErrInvalidQuery
	}

	return s.repo.GetTeachers(ctx, groupID, semesterID)
}

// Группы преподавателя по предмету
func (s *Service) GetGroups(
	ctx context.Context,
	teacherID, subjectID, semesterID int64,
) ([]GroupResult, error) {
	if !validIDs(teacherID, subjectID, semesterID) {
		return nil, ErrInvalidQuery
	}

	return s.repo.GetGroups(ctx, teacherID, subjectID, semesterID)
}

// Расписание группы на день
func (s *Service) GetSchedule(
	ctx context.Context,
	groupID, semesterID int64,
	weekday int,
) ([]Lesson, error) {
	if !validIDs(groupID, semesterID) || !validWeekday(weekday) {
		return nil, ErrInvalidQuery
	}

	return s.repo.GetSchedule(ctx, groupID, semesterID, weekday)
}

// Количество студентов по курсам в кабинете
func (s *Service) GetStudentCounts(
	ctx context.Context,
	classroomID, semesterID int64,
) ([]CourseCount, error) {
	if !validIDs(classroomID, semesterID) {
		return nil, ErrInvalidQuery
	}

	return s.repo.GetStudentCounts(ctx, classroomID, semesterID)
}
