package teaching

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type TeachingHandler struct {
	service *TeachingService
}

func NewTeachingHandler(service *TeachingService) *TeachingHandler {
	return &TeachingHandler{service: service}
}

type teachingRequest struct {
	TeacherID  int64 `json:"teacher_id"`
	SubjectID  int64 `json:"subject_id"`
	GroupID    int64 `json:"group_id"`
	SemesterID int64 `json:"semester_id"`
}

func (req teachingRequest) ToModel() TeachingAssignment {
	return TeachingAssignment{
		TeacherID:  req.TeacherID,
		SubjectID:  req.SubjectID,
		GroupID:    req.GroupID,
		SemesterID: req.SemesterID,
	}
}

func (h *TeachingHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/teachings", h.GetAll)
	mux.HandleFunc("GET /api/teachings/{id}", h.GetByID)
	mux.HandleFunc("POST /api/teachings", h.Create)
	mux.HandleFunc("PUT /api/teachings/{id}", h.Update)
	mux.HandleFunc("DELETE /api/teachings/{id}", h.Delete)
}

func (h *TeachingHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	assignments, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeaching)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, assignments)
}

func (h *TeachingHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	assignment, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeaching)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, assignment)
}

func (h *TeachingHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req teachingRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	assignment := req.ToModel()

	if err := h.service.Create(r.Context(), &assignment); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeaching)
		return
	}

	w.Header().Set("Location", "/api/teachings/"+strconv.FormatInt(assignment.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, assignment)
}

func (h *TeachingHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req teachingRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	assignment := req.ToModel()
	assignment.ID = id

	if err := h.service.Update(r.Context(), assignment); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeaching)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *TeachingHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidTeaching)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
