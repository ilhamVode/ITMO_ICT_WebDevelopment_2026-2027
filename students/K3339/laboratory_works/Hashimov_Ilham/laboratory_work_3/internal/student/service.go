package student

import (
	"context"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidStudent = errors.New("некорректные данные студента")

type StudentService struct {
	repo *StudentRepository
}

func NewStudentService(repo *StudentRepository) *StudentService {
	return &StudentService{repo: repo}
}

// Проверка данных студента
func validateStudent(student *Student) error {
	if student == nil {
		return fmt.Errorf("%w: студент не указан", ErrInvalidStudent)
	}

	student.FirstName = strings.TrimSpace(student.FirstName)
	student.LastName = strings.TrimSpace(student.LastName)

	if student.FirstName == "" {
		return fmt.Errorf("%w: имя не может быть пустым", ErrInvalidStudent)
	}

	if student.LastName == "" {
		return fmt.Errorf("%w: фамилия не может быть пустой", ErrInvalidStudent)
	}

	if student.GroupID <= 0 {
		return fmt.Errorf("%w: неверный ID группы", ErrInvalidStudent)
	}

	return nil
}

// Создание студента
func (s *StudentService) Create(ctx context.Context, student *Student) error {
	if err := validateStudent(student); err != nil {
		return err
	}

	student.ID = 0
	student.IsActive = true

	return s.repo.Create(ctx, student)
}

// Получение студента по ID
func (s *StudentService) GetByID(ctx context.Context, id int64) (Student, error) {
	if id <= 0 {
		return Student{}, ErrInvalidStudent
	}

	return s.repo.GetByID(ctx, id)
}

// Получение всех студентов
func (s *StudentService) GetAll(ctx context.Context) ([]Student, error) {
	return s.repo.GetAll(ctx)
}

// Обновление студента
func (s *StudentService) Update(ctx context.Context, student Student) error {
	if student.ID <= 0 {
		return ErrInvalidStudent
	}

	if err := validateStudent(&student); err != nil {
		return err
	}

	return s.repo.Update(ctx, student)
}

// Удаление студента
func (s *StudentService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidStudent
	}

	return s.repo.Delete(ctx, id)
}

// Отчисление сохраняет студента и его оценки в истории.
func (s *StudentService) Expel(ctx context.Context, id int64) error {
	student, err := s.GetByID(ctx, id)
	if err != nil {
		return err
	}
	student.IsActive = false
	return s.repo.Update(ctx, student)
}
