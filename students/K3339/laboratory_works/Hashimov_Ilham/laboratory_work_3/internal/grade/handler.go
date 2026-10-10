package grade

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type GradeHandler struct {
	service *GradeService
}

func NewGradeHandler(service *GradeService) *GradeHandler {
	return &GradeHandler{service: service}
}

type gradeRequest struct {
	StudentID  int64 `json:"student_id"`
	SubjectID  int64 `json:"subject_id"`
	SemesterID int64 `json:"semester_id"`
	Value      int   `json:"value"`
}

func (req gradeRequest) ToModel() Grade {
	return Grade{
		StudentID:  req.StudentID,
		SubjectID:  req.SubjectID,
		SemesterID: req.SemesterID,
		Value:      req.Value,
	}
}

func (h *GradeHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/grades", h.GetAll)
	mux.HandleFunc("GET /api/grades/{id}", h.GetByID)
	mux.HandleFunc("POST /api/grades", h.Create)
	mux.HandleFunc("PUT /api/grades/{id}", h.Update)
	mux.HandleFunc("DELETE /api/grades/{id}", h.Delete)
}

func (h *GradeHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	grades, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidGrade)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, grades)
}

func (h *GradeHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	grade, err := h.service.GetWithDetails(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidGrade)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, grade)
}

func (h *GradeHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req gradeRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	grade := req.ToModel()

	if err := h.service.Create(r.Context(), &grade); err != nil {
		httpapi.HandleError(w, err, ErrInvalidGrade)
		return
	}

	w.Header().Set("Location", "/api/grades/"+strconv.FormatInt(grade.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, grade)
}

func (h *GradeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req gradeRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	grade := req.ToModel()
	grade.ID = id

	if err := h.service.Update(r.Context(), grade); err != nil {
		httpapi.HandleError(w, err, ErrInvalidGrade)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *GradeHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidGrade)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
