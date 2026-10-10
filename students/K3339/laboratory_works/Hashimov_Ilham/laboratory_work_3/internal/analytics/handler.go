package analytics

import (
	"fmt"
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

// Читаем положительные целые числа из query-параметров
func readParams(r *http.Request, keys ...string) (map[string]int64, error) {
	result := make(map[string]int64, len(keys))

	for _, key := range keys {
		value, err := strconv.ParseInt(r.URL.Query().Get(key), 10, 64)

		if err != nil || value <= 0 {
			return nil, fmt.Errorf("некорректный параметр: %s", key)
		}

		result[key] = value
	}

	return result, nil
}

func analyticsResponse(w http.ResponseWriter, data any, err error) {
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidQuery)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, data)
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/analytics/lesson", h.GetLesson)
	mux.HandleFunc("GET /api/analytics/teachers", h.GetTeachers)
	mux.HandleFunc("GET /api/analytics/groups", h.GetGroups)
	mux.HandleFunc("GET /api/analytics/schedule", h.GetSchedule)
	mux.HandleFunc("GET /api/analytics/student-counts", h.GetStudentCounts)
}

// 1. Предмет на указанной паре
func (h *Handler) GetLesson(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(r,
		"group_id", "semester_id", "weekday", "lesson_number")

	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	// Сначала проверяем диапазон, затем преобразуем int64 в int
	if p["weekday"] > 7 || p["lesson_number"] > 2147483647 {
		httpapi.WriteError(w, http.StatusBadRequest, "некорректный день или номер пары")
		return
	}

	result, err := h.service.GetLesson(
		r.Context(),
		p["group_id"],
		p["semester_id"],
		int(p["weekday"]),
		int(p["lesson_number"]),
	)

	analyticsResponse(w, result, err)
}

// 2. Преподаватели группы
func (h *Handler) GetTeachers(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(r, "group_id", "semester_id")
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetTeachers(
		r.Context(),
		p["group_id"],
		p["semester_id"],
	)

	analyticsResponse(w, result, err)
}

// 3. Группы преподавателя по дисциплине
func (h *Handler) GetGroups(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(r,
		"teacher_id", "subject_id", "semester_id")

	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetGroups(
		r.Context(),
		p["teacher_id"],
		p["subject_id"],
		p["semester_id"],
	)

	analyticsResponse(w, result, err)
}

// 4. Расписание группы на день
func (h *Handler) GetSchedule(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(r,
		"group_id", "semester_id", "weekday")

	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if p["weekday"] > 7 {
		httpapi.WriteError(w, http.StatusBadRequest, "некорректный день недели")
		return
	}

	result, err := h.service.GetSchedule(
		r.Context(),
		p["group_id"],
		p["semester_id"],
		int(p["weekday"]),
	)

	analyticsResponse(w, result, err)
}

// 5. Количество студентов по курсам в кабинете
func (h *Handler) GetStudentCounts(w http.ResponseWriter, r *http.Request) {
	p, err := readParams(r, "classroom_id", "semester_id")
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.GetStudentCounts(
		r.Context(),
		p["classroom_id"],
		p["semester_id"],
	)

	analyticsResponse(w, result, err)
}
