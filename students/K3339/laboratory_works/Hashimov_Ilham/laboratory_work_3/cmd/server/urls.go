package main

import (
	"encoding/json"
	"net/http"

	"laboratory_work_3/internal/analytics"
	"laboratory_work_3/internal/auth"
	"laboratory_work_3/internal/classroom"
	"laboratory_work_3/internal/grade"
	"laboratory_work_3/internal/group"
	"laboratory_work_3/internal/report"
	"laboratory_work_3/internal/schedule"
	"laboratory_work_3/internal/semester"
	"laboratory_work_3/internal/student"
	"laboratory_work_3/internal/subject"
	"laboratory_work_3/internal/teacher"
	"laboratory_work_3/internal/teaching"
)

type Handlers struct {
	Student   *student.StudentHandler
	Group     *group.GroupHandler
	Teacher   *teacher.TeacherHandler
	Classroom *classroom.ClassroomHandler
	Subject   *subject.SubjectHandler
	Semester  *semester.SemesterHandler
	Teaching  *teaching.TeachingHandler
	Schedule  *schedule.ScheduleHandler
	Grade     *grade.GradeHandler
	Analytics *analytics.Handler
	Report    *report.Handler
	Auth      *auth.Handler
}

func newRouter(h Handlers) *http.ServeMux {
	mux := http.NewServeMux()

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")

		_ = json.NewEncoder(w).Encode(map[string]string{
			"status": "ok",
		})
	})

	api := http.NewServeMux()
	h.Student.RegisterRoutes(api)
	h.Group.RegisterRoutes(api)
	h.Teacher.RegisterRoutes(api)
	h.Classroom.RegisterRoutes(api)
	h.Subject.RegisterRoutes(api)
	h.Semester.RegisterRoutes(api)
	h.Teaching.RegisterRoutes(api)
	h.Schedule.RegisterRoutes(api)
	h.Grade.RegisterRoutes(api)
	h.Analytics.RegisterRoutes(api)
	h.Report.RegisterRoutes(api)
	h.Auth.RegisterRoutes(mux)
	mux.Handle("/api/", h.Auth.RequireAuth(api))

	return mux
}
