package group

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type GroupHandler struct {
	service *GroupService
}

func NewGroupHandler(service *GroupService) *GroupHandler {
	return &GroupHandler{service: service}
}

// Данные, которые присылает клиент
type groupRequest struct {
	Name   string `json:"name"`
	Course int    `json:"course"`
}

// Ответ с вложенными студентами
type groupStudentResponse struct {
	ID        int64   `json:"id"`
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	BirthDate *string `json:"birth_date"`
	GroupID   int64   `json:"group_id"`
	IsActive  bool    `json:"is_active"`
}

type groupDetailsResponse struct {
	ID       int64                  `json:"id"`
	Name     string                 `json:"name"`
	Course   int                    `json:"course"`
	Students []groupStudentResponse `json:"students"`
}

// Подготавливаем группу со студентами для JSON
func toGroupDetails(group StudyGroup) groupDetailsResponse {
	result := groupDetailsResponse{
		ID:       group.ID,
		Name:     group.Name,
		Course:   group.Course,
		Students: make([]groupStudentResponse, 0, len(group.Students)),
	}

	for _, student := range group.Students {
		var birthDate *string

		if student.BirthDate != nil {
			date := student.BirthDate.Format("2006-01-02")
			birthDate = &date
		}

		result.Students = append(result.Students, groupStudentResponse{
			ID:        student.ID,
			FirstName: student.FirstName,
			LastName:  student.LastName,
			BirthDate: birthDate,
			GroupID:   student.GroupID,
			IsActive:  student.IsActive,
		})
	}

	return result
}

// Регистрируем маршруты
func (h *GroupHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/groups", h.GetAll)
	mux.HandleFunc("GET /api/groups/{id}", h.GetByID)
	mux.HandleFunc("POST /api/groups", h.Create)
	mux.HandleFunc("PUT /api/groups/{id}", h.Update)
	mux.HandleFunc("DELETE /api/groups/{id}", h.Delete)
}

// GET /api/groups
func (h *GroupHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	groups, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidGroup)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, groups)
}

// GET /api/groups/{id}
func (h *GroupHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	group, err := h.service.GetWithStudents(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidGroup)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, toGroupDetails(group))
}

// POST /api/groups
func (h *GroupHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req groupRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	group := StudyGroup{
		Name:   req.Name,
		Course: req.Course,
	}

	if err := h.service.Create(r.Context(), &group); err != nil {
		httpapi.HandleError(w, err, ErrInvalidGroup)
		return
	}

	w.Header().Set("Location", "/api/groups/"+strconv.FormatInt(group.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, group)
}

// PUT /api/groups/{id}
func (h *GroupHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req groupRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	group := StudyGroup{
		ID:     id,
		Name:   req.Name,
		Course: req.Course,
	}

	if err := h.service.Update(r.Context(), group); err != nil {
		httpapi.HandleError(w, err, ErrInvalidGroup)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// DELETE /api/groups/{id}
func (h *GroupHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidGroup)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
