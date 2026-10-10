package student

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

type StudentHandler struct {
	service *StudentService
}

func NewStudentHandler(service *StudentService) *StudentHandler {
	return &StudentHandler{service: service}
}

// Данные, которые клиент отправляет серверу
type studentRequest struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	GroupID   int64   `json:"group_id"`
}

type updateStudentRequest struct {
	studentRequest
	IsActive *bool `json:"is_active"`
}

// Преобразуем запрос в модель студента
func (req studentRequest) ToModel() (Student, error) {
	student := Student{
		FirstName: req.FirstName,
		LastName:  req.LastName,
		GroupID:   req.GroupID,
	}

	if req.BirthDate != nil {
		date, err := time.Parse("2006-01-02", *req.BirthDate)
		if err != nil {
			return Student{}, errors.New("birth_date должен иметь формат YYYY-MM-DD")
		}

		student.BirthDate = &date
	}

	return student, nil
}

// Преобразуем модель в JSON-ответ
func studentResponse(student Student) map[string]any {
	var birthDate any

	if student.BirthDate != nil {
		birthDate = student.BirthDate.Format("2006-01-02")
	}

	return map[string]any{
		"id":         student.ID,
		"first_name": student.FirstName,
		"last_name":  student.LastName,
		"birth_date": birthDate,
		"group_id":   student.GroupID,
		"is_active":  student.IsActive,
	}
}

// Отправка JSON
func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(data)
}

func writeError(w http.ResponseWriter, status int, message string) {
	writeJSON(w, status, map[string]string{"error": message})
}

// Чтение JSON из тела запроса
func readJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("ожидался один JSON-объект")
	}

	return nil
}

// Получение ID из URL
func getID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("некорректный ID")
	}

	return id, nil
}

// Преобразуем ошибки в HTTP-ответы
func handleStudentError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, ErrInvalidStudent):
		writeError(w, http.StatusBadRequest, err.Error())

	case errors.Is(err, gorm.ErrRecordNotFound):
		writeError(w, http.StatusNotFound, "студент не найден")

	default:
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case "23503":
				writeError(w, http.StatusBadRequest, "нарушена связь с другой таблицей")
				return
			case "23505":
				writeError(w, http.StatusConflict, "запись уже существует")
				return
			case "23514", "23502":
				writeError(w, http.StatusBadRequest, "нарушены ограничения базы данных")
				return
			}
		}

		log.Printf("Ошибка Student API: %v", err)
		writeError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
	}
}

// Регистрация маршрутов
func (h *StudentHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/students", h.GetAll)
	mux.HandleFunc("GET /api/students/{id}", h.GetByID)
	mux.HandleFunc("POST /api/students", h.Create)
	mux.HandleFunc("PUT /api/students/{id}", h.Update)
	mux.HandleFunc("DELETE /api/students/{id}", h.Delete)
	mux.HandleFunc("POST /api/students/{id}/expel", h.Expel)
}

// GET /api/students
func (h *StudentHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	students, err := h.service.GetAll(r.Context())
	if err != nil {
		handleStudentError(w, err)
		return
	}

	result := make([]map[string]any, 0, len(students))

	for _, student := range students {
		result = append(result, studentResponse(student))
	}

	writeJSON(w, http.StatusOK, result)
}

// GET /api/students/{id}
func (h *StudentHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	student, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		handleStudentError(w, err)
		return
	}

	writeJSON(w, http.StatusOK, studentResponse(student))
}

// POST /api/students
func (h *StudentHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req studentRequest

	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON: "+err.Error())
		return
	}

	student, err := req.ToModel()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Create(r.Context(), &student); err != nil {
		handleStudentError(w, err)
		return
	}

	w.Header().Set("Location", "/api/students/"+strconv.FormatInt(student.ID, 10))
	writeJSON(w, http.StatusCreated, studentResponse(student))
}

// PUT /api/students/{id}
func (h *StudentHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req updateStudentRequest

	if err := readJSON(w, r, &req); err != nil {
		writeError(w, http.StatusBadRequest, "некорректный JSON: "+err.Error())
		return
	}

	if req.IsActive == nil {
		writeError(w, http.StatusBadRequest, "необходимо указать is_active")
		return
	}

	student, err := req.studentRequest.ToModel()
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	student.ID = id
	student.IsActive = *req.IsActive

	if err := h.service.Update(r.Context(), student); err != nil {
		handleStudentError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/students/{id}
func (h *StudentHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		handleStudentError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *StudentHandler) Expel(w http.ResponseWriter, r *http.Request) {
	id, err := getID(r)
	if err != nil {
		writeError(w, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.service.Expel(r.Context(), id); err != nil {
		handleStudentError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
