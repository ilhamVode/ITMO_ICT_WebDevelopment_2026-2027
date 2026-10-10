package semester

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type SemesterHandler struct {
	service *SemesterService
}

func NewSemesterHandler(service *SemesterService) *SemesterHandler {
	return &SemesterHandler{service: service}
}

type semesterRequest struct {
	AcademicYear string `json:"academic_year"`
	Number       int    `json:"number"`
}

func (h *SemesterHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/semesters", h.GetAll)
	mux.HandleFunc("GET /api/semesters/{id}", h.GetByID)
	mux.HandleFunc("POST /api/semesters", h.Create)
	mux.HandleFunc("PUT /api/semesters/{id}", h.Update)
	mux.HandleFunc("DELETE /api/semesters/{id}", h.Delete)
}

func (h *SemesterHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	semesters, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidSemester)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, semesters)
}

func (h *SemesterHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	semester, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidSemester)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, semester)
}

func (h *SemesterHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req semesterRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	semester := Semester{
		AcademicYear: req.AcademicYear,
		Number:       req.Number,
	}

	if err := h.service.Create(r.Context(), &semester); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSemester)
		return
	}

	w.Header().Set("Location", "/api/semesters/"+strconv.FormatInt(semester.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, semester)
}

func (h *SemesterHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req semesterRequest
	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	semester := Semester{
		ID:           id,
		AcademicYear: req.AcademicYear,
		Number:       req.Number,
	}

	if err := h.service.Update(r.Context(), semester); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSemester)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SemesterHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSemester)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
