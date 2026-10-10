package semester

import (
	"context"
	"errors"
	"strconv"
	"strings"
)

var ErrInvalidSemester = errors.New("некорректные данные: semester")

type SemesterService struct {
	repo *SemesterRepository
}

func NewSemesterService(repo *SemesterRepository) *SemesterService {
	return &SemesterService{repo: repo}
}

func validateSemester(semester *Semester) error {
	if semester == nil {
		return ErrInvalidSemester
	}

	semester.AcademicYear = strings.TrimSpace(semester.AcademicYear)
	year := semester.AcademicYear

	if len(year) != 9 || year[4] != '-' ||
		(semester.Number != 1 && semester.Number != 2) {
		return ErrInvalidSemester
	}

	start, err1 := strconv.Atoi(year[:4])
	end, err2 := strconv.Atoi(year[5:])

	if err1 != nil || err2 != nil || start <= 0 || end != start+1 {
		return ErrInvalidSemester
	}

	return nil
}

func (s *SemesterService) Create(ctx context.Context, semester *Semester) error {
	if err := validateSemester(semester); err != nil {
		return err
	}

	semester.ID = 0
	return s.repo.Create(ctx, semester)
}

func (s *SemesterService) GetByID(ctx context.Context, id int64) (Semester, error) {
	if id <= 0 {
		return Semester{}, ErrInvalidSemester
	}

	return s.repo.GetByID(ctx, id)
}

func (s *SemesterService) GetAll(ctx context.Context) ([]Semester, error) {
	return s.repo.GetAll(ctx)
}

func (s *SemesterService) Update(ctx context.Context, semester Semester) error {
	if semester.ID <= 0 {
		return ErrInvalidSemester
	}

	if err := validateSemester(&semester); err != nil {
		return err
	}

	return s.repo.Update(ctx, semester)
}

func (s *SemesterService) Delete(ctx context.Context, id int64) error {
	if id <= 0 {
		return ErrInvalidSemester
	}

	return s.repo.Delete(ctx, id)
}
