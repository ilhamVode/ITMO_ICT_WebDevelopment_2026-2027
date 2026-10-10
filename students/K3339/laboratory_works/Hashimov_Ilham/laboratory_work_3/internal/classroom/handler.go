package classroom

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type ClassroomHandler struct {
	service *ClassroomService
}

func NewClassroomHandler(service *ClassroomService) *ClassroomHandler {
	return &ClassroomHandler{service: service}
}

type classroomRequest struct {
	Number   string `json:"number"`
	Building string `json:"building"`
}

func (h *ClassroomHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/classrooms", h.GetAll)
	mux.HandleFunc("GET /api/classrooms/{id}", h.GetByID)
	mux.HandleFunc("POST /api/classrooms", h.Create)
	mux.HandleFunc("PUT /api/classrooms/{id}", h.Update)
	mux.HandleFunc("DELETE /api/classrooms/{id}", h.Delete)
}

func (h *ClassroomHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	classrooms, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidClassroom)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, classrooms)
}

func (h *ClassroomHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	classroom, err := h.service.GetByID(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidClassroom)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, classroom)
}

func (h *ClassroomHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req classroomRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	classroom := Classroom{
		Number:   req.Number,
		Building: req.Building,
	}

	if err := h.service.Create(r.Context(), &classroom); err != nil {
		httpapi.HandleError(w, err, ErrInvalidClassroom)
		return
	}

	w.Header().Set("Location", "/api/classrooms/"+strconv.FormatInt(classroom.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, classroom)
}

func (h *ClassroomHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req classroomRequest
	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	classroom := Classroom{
		ID:       id,
		Number:   req.Number,
		Building: req.Building,
	}

	if err := h.service.Update(r.Context(), classroom); err != nil {
		httpapi.HandleError(w, err, ErrInvalidClassroom)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ClassroomHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidClassroom)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
