package schedule

import (
	"net/http"
	"strconv"

	"laboratory_work_3/internal/httpapi"
)

type ScheduleHandler struct {
	service *ScheduleService
}

func NewScheduleHandler(service *ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{service: service}
}

type scheduleRequest struct {
	AssignmentID int64 `json:"assignment_id"`
	ClassroomID  int64 `json:"classroom_id"`
	Weekday      int   `json:"weekday"`
	LessonNumber int   `json:"lesson_number"`
}

func (req scheduleRequest) ToModel() Schedule {
	return Schedule{
		AssignmentID: req.AssignmentID,
		ClassroomID:  req.ClassroomID,
		Weekday:      req.Weekday,
		LessonNumber: req.LessonNumber,
	}
}

func (h *ScheduleHandler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/schedules", h.GetAll)
	mux.HandleFunc("GET /api/schedules/{id}", h.GetByID)
	mux.HandleFunc("POST /api/schedules", h.Create)
	mux.HandleFunc("PUT /api/schedules/{id}", h.Update)
	mux.HandleFunc("DELETE /api/schedules/{id}", h.Delete)
}

func (h *ScheduleHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	schedules, err := h.service.GetAll(r.Context())
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidSchedule)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, schedules)
}

func (h *ScheduleHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	schedule, err := h.service.GetWithDetails(r.Context(), id)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidSchedule)
		return
	}

	httpapi.WriteJSON(w, http.StatusOK, schedule)
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req scheduleRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	schedule := req.ToModel()

	if err := h.service.Create(r.Context(), &schedule); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSchedule)
		return
	}

	w.Header().Set("Location", "/api/schedules/"+strconv.FormatInt(schedule.ID, 10))
	httpapi.WriteJSON(w, http.StatusCreated, schedule)
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	var req scheduleRequest

	if err := httpapi.ReadJSON(w, r, &req); err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	schedule := req.ToModel()
	schedule.ID = id

	if err := h.service.Update(r.Context(), schedule); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSchedule)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *ScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, err := httpapi.PathID(r)
	if err != nil {
		httpapi.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.service.Delete(r.Context(), id); err != nil {
		httpapi.HandleError(w, err, ErrInvalidSchedule)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
