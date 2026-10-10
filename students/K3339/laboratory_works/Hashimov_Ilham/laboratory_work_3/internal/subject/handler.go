package subject

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type SubjectHandler struct {
	service *SubjectService
}

func NewSubjectHandler(service *SubjectService) *SubjectHandler {
	return &SubjectHandler{service: service}
}

type subjectRequest struct {
	Name string `json:"name"`
}

func (h *SubjectHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/subjects", h.GetAll)
	mux.HandleFunc("GET /api/subjects/{id}", h.GetByID)
	mux.HandleFunc("POST /api/subjects", h.Create)
	mux.HandleFunc("PUT /api/subjects/{id}", h.Update)
	mux.HandleFunc("DELETE /api/subjects/{id}", h.Delete)
}

func (h *SubjectHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	subjects, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidSubject)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, subjects)
}

func (h *SubjectHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	subject, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidSubject)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, subject)
}

func (h *SubjectHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req subjectRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	subject := Subject{Name: req.Name}

	if err := h.service.Create(r.Context(), &subject); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSubject)
		return
	}

	w.Header().Set("Location", "/api/subjects/"+strconv.FormatInt(subject.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, subject)
}

func (h *SubjectHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req subjectRequest
	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	subject := Subject{
		ID:   id,
		Name: req.Name,
	}

	if err := h.service.Update(r.Context(), subject); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSubject)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *SubjectHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSubject)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
