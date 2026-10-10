package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"laboratory_work_3/internal/analytics"
	"laboratory_work_3/internal/auth"
	"laboratory_work_3/internal/classroom"
	"laboratory_work_3/internal/database"
	"laboratory_work_3/internal/grade"
	"laboratory_work_3/internal/group"
	"laboratory_work_3/internal/report"
	"laboratory_work_3/internal/schedule"
	"laboratory_work_3/internal/semester"
	"laboratory_work_3/internal/student"
	"laboratory_work_3/internal/subject"
	"laboratory_work_3/internal/teacher"
	"laboratory_work_3/internal/teaching"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil && !os.IsNotExist(err) {
		log.Fatal("Ошибка загрузки .env: ", err)
	}

	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("Не задан DATABASE_DSN")
	}

	db, err := database.NewPostgres(dsn)
	if err != nil {
		log.Fatal(err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatal(err)
	}
	defer sqlDB.Close()

	// Инициализируем схему базы данных
	if err := database.InitSchema(db); err != nil {
		log.Fatal(err)
	}

	// Student
	studentRepo := student.NewStudentRepository(db)
	studentService := student.NewStudentService(studentRepo)
	studentHandler := student.NewStudentHandler(studentService)

	// Group
	groupRepo := group.NewGroupRepository(db)
	groupService := group.NewGroupService(groupRepo)
	groupHandler := group.NewGroupHandler(groupService)

	// Teacher
	teacherRepo := teacher.NewTeacherRepository(db)
	teacherService := teacher.NewTeacherService(teacherRepo)
	teacherHandler := teacher.NewTeacherHandler(teacherService)

	// Classroom
	classroomRepo := classroom.NewClassroomRepository(db)
	classroomService := classroom.NewClassroomService(classroomRepo)
	classroomHandler := classroom.NewClassroomHandler(classroomService)

	// Subject
	subjectRepo := subject.NewSubjectRepository(db)
	subjectService := subject.NewSubjectService(subjectRepo)
	subjectHandler := subject.NewSubjectHandler(subjectService)

	// Semester
	semesterRepo := semester.NewSemesterRepository(db)
	semesterService := semester.NewSemesterService(semesterRepo)
	semesterHandler := semester.NewSemesterHandler(semesterService)

	// Teaching
	teachingRepo := teaching.NewTeachingRepository(db)
	teachingService := teaching.NewTeachingService(teachingRepo)
	teachingHandler := teaching.NewTeachingHandler(teachingService)

	// Schedule
	scheduleRepo := schedule.NewScheduleRepository(db)
	scheduleService := schedule.NewScheduleService(scheduleRepo)
	scheduleHandler := schedule.NewScheduleHandler(scheduleService)

	// Grade
	gradeRepo := grade.NewGradeRepository(db)
	gradeService := grade.NewGradeService(gradeRepo)
	gradeHandler := grade.NewGradeHandler(gradeService)

	// Analytics
	analyticsRepo := analytics.NewRepository(db)
	analyticsService := analytics.NewService(analyticsRepo)
	analyticsHandler := analytics.NewHandler(analyticsService)

	// Report
	reportRepo := report.NewRepository(db)
	reportService := report.NewService(reportRepo)
	reportHandler := report.NewHandler(reportService)

	// Auth
	authRepo := auth.NewRepository(db)
	authService := auth.NewService(authRepo)
	authHandler := auth.NewHandler(authService)

	// Собираем обработчики
	handlers := Handlers{
		Student:   studentHandler,
		Group:     groupHandler,
		Teacher:   teacherHandler,
		Classroom: classroomHandler,
		Subject:   subjectHandler,
		Semester:  semesterHandler,
		Teaching:  teachingHandler,
		Schedule:  scheduleHandler,
		Grade:     gradeHandler,
		Analytics: analyticsHandler,
		Report:    reportHandler,
		Auth:      authHandler,
	}

	mux := newRouter(handlers)
	port := os.Getenv("SERVER_ADDR")
	if port == "" {
		port = ":8080"
	}
	log.Printf("Сервер запускается на %s", port)
	server := &http.Server{
		Addr: port, Handler: mux,
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       60 * time.Second,
	}
	log.Fatal(server.ListenAndServe())
}
