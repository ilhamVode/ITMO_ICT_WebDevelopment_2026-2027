package teacher

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type TeacherHandler struct {
	service *TeacherService
}

func NewTeacherHandler(service *TeacherService) *TeacherHandler {
	return &TeacherHandler{service: service}
}

type teacherRequest struct {
	FirstName        string  `json:"first_name"`
	LastName         string  `json:"last_name"`
	MiddleName       *string `json:"middle_name"`
	FixedClassroomID *int64  `json:"fixed_classroom_id"`
}

type teacherUpdateRequest struct {
	teacherRequest
	IsActive *bool `json:"is_active"`
}

func (req teacherRequest) ToModel() Teacher {
	return Teacher{
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		MiddleName:       req.MiddleName,
		FixedClassroomID: req.FixedClassroomID,
	}
}

func (h *TeacherHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/teachers", h.GetAll)
	mux.HandleFunc("GET /api/teachers/{id}", h.GetByID)
	mux.HandleFunc("POST /api/teachers", h.Create)
	mux.HandleFunc("PUT /api/teachers/{id}", h.Update)
	mux.HandleFunc("DELETE /api/teachers/{id}", h.Delete)
	mux.HandleFunc("POST /api/teachers/{id}/dismiss", h.Dismiss)
}

func (h *TeacherHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	teachers, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeacher)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, teachers)
}

func (h *TeacherHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	teacher, err := h.service.GetWithSubjects(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeacher)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, teacher)
}

func (h *TeacherHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req teacherRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	teacher := req.ToModel()

	if err := h.service.Create(r.Context(), &teacher); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeacher)
		return
	}

	w.Header().Set("Location", "/api/teachers/"+strconv.FormatInt(teacher.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, teacher)
}

func (h *TeacherHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req teacherUpdateRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if req.IsActive == nil {
		httpapi.WriteError(w, http.StatusBadRequest, "необходимо указать is_active")
		return
	}

	teacher := req.teacherRequest.ToModel()
	teacher.ID = id
	teacher.IsActive = *req.IsActive

	if err := h.service.Update(r.Context(), teacher); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeacher)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeacherHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeacher)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeacherHandler) Dismiss(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Dismiss(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeacher)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
