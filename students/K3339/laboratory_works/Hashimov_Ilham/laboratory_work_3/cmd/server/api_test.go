package main

import (
	"bytes"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"
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

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TEST_DATABASE_DSN задаётся явно. Существующие таблицы не изменяются:
// каждый запуск создаёт свою схему и удаляет только её после проверки.
func testDatabase(t *testing.T) *gorm.DB {
	t.Helper()
	dsn := os.Getenv("TEST_DATABASE_DSN")
	if dsn == "" {
		t.Skip("для интеграционного теста нужен TEST_DATABASE_DSN")
	}
	admin, err := database.NewPostgres(dsn)
	if err != nil {
		t.Fatal("не удалось подключиться к тестовой PostgreSQL")
	}
	schema := fmt.Sprintf("lab3_test_%d", time.Now().UnixNano())
	if err := admin.Exec("CREATE SCHEMA " + schema).Error; err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() {
		if err := admin.Exec("DROP SCHEMA " + schema + " CASCADE").Error; err != nil {
			t.Error(err)
		}
		sqlDB, _ := admin.DB()
		_ = sqlDB.Close()
	})
	config, err := pgx.ParseConfig(dsn)
	if err != nil {
		t.Fatal("не удалось прочитать TEST_DATABASE_DSN")
	}
	config.RuntimeParams["search_path"] = schema
	config.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol
	sqlDB := stdlib.OpenDB(*config)
	db, err := gorm.Open(postgres.New(postgres.Config{Conn: sqlDB}), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = sqlDB.Close() })
	// init.sql применяется относительно корня лабораторной.
	t.Chdir(filepath.Join("..", ".."))
	if err := database.InitSchema(db); err != nil {
		t.Fatal(err)
	}
	// Повторная инициализация должна сохранять данные и не давать ошибок.
	if err := database.InitSchema(db); err != nil {
		t.Fatal(err)
	}
	return db
}

func testRouter(db *gorm.DB) http.Handler {
	return newRouter(Handlers{
		Student:   student.NewStudentHandler(student.NewStudentService(student.NewStudentRepository(db))),
		Group:     group.NewGroupHandler(group.NewGroupService(group.NewGroupRepository(db))),
		Teacher:   teacher.NewTeacherHandler(teacher.NewTeacherService(teacher.NewTeacherRepository(db))),
		Classroom: classroom.NewClassroomHandler(classroom.NewClassroomService(classroom.NewClassroomRepository(db))),
		Subject:   subject.NewSubjectHandler(subject.NewSubjectService(subject.NewSubjectRepository(db))),
		Semester:  semester.NewSemesterHandler(semester.NewSemesterService(semester.NewSemesterRepository(db))),
		Teaching:  teaching.NewTeachingHandler(teaching.NewTeachingService(teaching.NewTeachingRepository(db))),
		Schedule:  schedule.NewScheduleHandler(schedule.NewScheduleService(schedule.NewScheduleRepository(db))),
		Grade:     grade.NewGradeHandler(grade.NewGradeService(grade.NewGradeRepository(db))),
		Analytics: analytics.NewHandler(analytics.NewService(analytics.NewRepository(db))),
		Report:    report.NewHandler(report.NewService(report.NewRepository(db))),
		Auth:      auth.NewHandler(auth.NewService(auth.NewRepository(db))),
	})
}

func apiRequest(t *testing.T, router http.Handler, method, url, token string, body any, status int) []byte {
	t.Helper()
	var input []byte
	if body != nil {
		var err error
		input, err = json.Marshal(body)
		if err != nil {
			t.Fatal(err)
		}
	}
	req := httptest.NewRequest(method, url, bytes.NewReader(input))
	if token != "" {
		req.Header.Set("Authorization", "Token "+token)
	}
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if w.Code != status {
		t.Fatalf("%s %s: статус %d, ожидался %d; %s", method, url, w.Code, status, w.Body.String())
	}
	return w.Body.Bytes()
}

func TestCollegeAPI(t *testing.T) {
	db := testDatabase(t)
	router := testRouter(db)
	apiRequest(t, router, "GET", "/health", "", nil, 200)
	apiRequest(t, router, "GET", "/api/groups", "", nil, 401)
	credentials := map[string]any{"username": "student_test", "password": "password123"}
	registered := apiRequest(t, router, "POST", "/api/auth/users", "", credentials, 201)
	if bytes.Contains(registered, []byte("password")) {
		t.Fatal("пароль попал в ответ")
	}
	apiRequest(t, router, "POST", "/api/auth/users", "", credentials, 409)
	apiRequest(t, router, "POST", "/api/auth/token/login", "", map[string]any{"username": "student_test", "password": "wrong"}, 401)
	var login struct {
		Token string `json:"auth_token"`
	}
	if err := json.Unmarshal(apiRequest(t, router, "POST", "/api/auth/token/login", "", credentials, 200), &login); err != nil {
		t.Fatal(err)
	}
	token := login.Token
	apiRequest(t, router, "GET", "/api/auth/users/me", token, nil, 200)
	create := func(resource string, body any) int64 {
		t.Helper()
		var item struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(apiRequest(t, router, "POST", "/api/"+resource, token, body, 201), &item); err != nil {
			t.Fatal(err)
		}
		return item.ID
	}
	g1 := create("groups", map[string]any{"name": "К3339", "course": 3})
	g2 := create("groups", map[string]any{"name": "К2339", "course": 2})
	g3 := create("groups", map[string]any{"name": "К1339", "course": 1})
	s1 := create("semesters", map[string]any{"academic_year": "2026-2027", "number": 1})
	s2 := create("semesters", map[string]any{"academic_year": "2026-2027", "number": 2})
	room1 := create("classrooms", map[string]any{"number": "301", "building": "Главный"})
	room2 := create("classrooms", map[string]any{"number": "302", "building": "Главный"})
	mathID := create("subjects", map[string]any{"name": "Математика"})
	physicsID := create("subjects", map[string]any{"name": "Физика"})
	t1 := create("teachers", map[string]any{"first_name": "Иван", "last_name": "Петров"})
	t2 := create("teachers", map[string]any{"first_name": "Мария", "last_name": "Соколова", "fixed_classroom_id": room1})
	st1 := create("students", map[string]any{"first_name": "Ильхам", "last_name": "Хашимов", "birth_date": "2005-01-01", "group_id": g1})
	st2 := create("students", map[string]any{"first_name": "Анна", "last_name": "Иванова", "group_id": g1})
	st3 := create("students", map[string]any{"first_name": "Олег", "last_name": "Иванов", "group_id": g3})
	assignment := func(teacherID, subjectID, groupID, semesterID int64) map[string]any {
		return map[string]any{"teacher_id": teacherID, "subject_id": subjectID, "group_id": groupID, "semester_id": semesterID}
	}
	a1 := create("teachings", assignment(t1, mathID, g1, s1))
	a2 := create("teachings", assignment(t2, mathID, g2, s1))
	a3 := create("teachings", assignment(t1, physicsID, g1, s1))
	a4 := create("teachings", assignment(t1, mathID, g3, s1))
	a5 := create("teachings", assignment(t1, mathID, g1, s2))
	apiRequest(t, router, "POST", "/api/teachings", token, assignment(t1, mathID, g1, s1), 409)
	slot := func(assignmentID, roomID int64, day, lesson int) map[string]any {
		return map[string]any{"assignment_id": assignmentID, "classroom_id": roomID, "weekday": day, "lesson_number": lesson}
	}
	create("schedules", slot(a1, room1, 1, 1))
	apiRequest(t, router, "POST", "/api/schedules", token, slot(a3, room2, 1, 1), 409) // группа
	apiRequest(t, router, "POST", "/api/schedules", token, slot(a4, room2, 1, 1), 409) // преподаватель
	apiRequest(t, router, "POST", "/api/schedules", token, slot(a2, room1, 1, 1), 409) // кабинет
	create("schedules", slot(a2, room2, 1, 1))
	create("schedules", slot(a5, room1, 1, 1)) // другой семестр
	create("schedules", slot(a3, room1, 1, 2))
	apiRequest(t, router, "PUT", fmt.Sprintf("/api/teachings/%d", a2), token, assignment(t2, mathID, g1, s1), 409)
	apiRequest(t, router, "GET", "/api/analytics/schedule?group_id=bad", token, nil, 400)
	urls := []string{
		fmt.Sprintf("lesson?group_id=%d&semester_id=%d&weekday=1&lesson_number=1", g1, s1),
		fmt.Sprintf("teachers?group_id=%d&semester_id=%d", g1, s1),
		fmt.Sprintf("groups?teacher_id=%d&subject_id=%d&semester_id=%d", t1, mathID, s1),
		fmt.Sprintf("schedule?group_id=%d&semester_id=%d&weekday=1", g1, s1),
		fmt.Sprintf("student-counts?classroom_id=%d&semester_id=%d", room1, s1),
	}
	for i, url := range urls {
		var rows []map[string]any
		if err := json.Unmarshal(apiRequest(t, router, "GET", "/api/analytics/"+url, token, nil, 200), &rows); err != nil {
			t.Fatal(err)
		}
		expected := []int{1, 1, 2, 2, 1}[i]
		if len(rows) != expected {
			t.Fatalf("аналитика %d: ожидалось %d строк, получено %d", i+1, expected, len(rows))
		}
		if i == 4 && rows[0]["student_count"] != float64(2) {
			t.Fatal("студенты посчитаны повторно")
		}
	}
	var groupDetails struct {
		Students []json.RawMessage `json:"students"`
	}
	_ = json.Unmarshal(apiRequest(t, router, "GET", fmt.Sprintf("/api/groups/%d", g1), token, nil, 200), &groupDetails)
	if len(groupDetails.Students) != 2 {
		t.Fatal("нет вложенных студентов")
	}
	var teacherDetails struct {
		Subjects []json.RawMessage `json:"subjects"`
	}
	_ = json.Unmarshal(apiRequest(t, router, "GET", fmt.Sprintf("/api/teachers/%d", t1), token, nil, 200), &teacherDetails)
	if len(teacherDetails.Subjects) != 2 {
		t.Fatal("неверный вложенный список m2m")
	}
	gradeData := func(studentID, subjectID, semesterID int64, value int) map[string]any {
		return map[string]any{"student_id": studentID, "subject_id": subjectID, "semester_id": semesterID, "value": value}
	}
	grade1 := create("grades", gradeData(st1, mathID, s1, 5))
	create("grades", gradeData(st1, physicsID, s1, 4))
	create("grades", gradeData(st2, mathID, s1, 3))
	create("grades", gradeData(st1, mathID, s2, 2))
	apiRequest(t, router, "POST", "/api/grades", token, gradeData(st3, physicsID, s1, 4), 400)
	apiRequest(t, router, "POST", "/api/grades", token, gradeData(st1, mathID, s1, 5), 409)
	reportURL := fmt.Sprintf("/api/reports/statement?group_id=%d&semester_id=%d", g1, s1)
	var statement report.Statement
	if err := json.Unmarshal(apiRequest(t, router, "GET", reportURL, token, nil, 200), &statement); err != nil {
		t.Fatal(err)
	}
	if statement.StudentCount != 2 || statement.GradeCount != 3 || statement.Average == nil || *statement.Average != 4 {
		t.Fatalf("неверная ведомость: %+v", statement)
	}
	csvData := apiRequest(t, router, "GET", reportURL+"&format=csv", token, nil, 200)
	reader := csv.NewReader(strings.NewReader(strings.TrimPrefix(string(csvData), "\uFEFF")))
	reader.Comma = ';'
	reader.FieldsPerRecord = -1
	rows, err := reader.ReadAll()
	if err != nil || len(rows) != 8 {
		t.Fatalf("неверная CSV-ведомость: %v", err)
	}
	apiRequest(t, router, "PUT", fmt.Sprintf("/api/grades/%d", grade1), token, gradeData(st1, mathID, s1, 2), 204)
	_ = json.Unmarshal(apiRequest(t, router, "GET", reportURL, token, nil, 200), &statement)
	if *statement.Average != 3 {
		t.Fatal("изменение оценки не попало в ведомость")
	}
	apiRequest(t, router, "POST", fmt.Sprintf("/api/students/%d/expel", st2), token, nil, 204)
	apiRequest(t, router, "PUT", fmt.Sprintf("/api/grades/%d", grade1), token, gradeData(st2, mathID, s1, 5), 400)
	apiRequest(t, router, "DELETE", fmt.Sprintf("/api/students/%d", st1), token, nil, 204)
	apiRequest(t, router, "GET", fmt.Sprintf("/api/grades/%d", grade1), token, nil, 404)
	apiRequest(t, router, "POST", fmt.Sprintf("/api/teachers/%d/dismiss", t2), token, nil, 204)
	apiRequest(t, router, "POST", "/api/teachings", token, assignment(t2, physicsID, g2, s1), 400)
	apiRequest(t, router, "DELETE", fmt.Sprintf("/api/teachers/%d", t2), token, nil, 204)
	apiRequest(t, router, "GET", fmt.Sprintf("/api/teachings/%d", a2), token, nil, 404)

	// Одновременные запросы не должны оба занять одну пару.
	server := httptest.NewServer(router)
	defer server.Close()
	statuses := make(chan int, 2)
	var wg sync.WaitGroup
	for _, a := range []int64{a1, a3} {
		wg.Add(1)
		go func(id int64) {
			defer wg.Done()
			input, _ := json.Marshal(slot(id, room1, 2, 1))
			req, _ := http.NewRequest("POST", server.URL+"/api/schedules", bytes.NewReader(input))
			req.Header.Set("Authorization", "Token "+token)
			response, err := server.Client().Do(req)
			if err != nil {
				statuses <- 0
				return
			}
			_, _ = io.Copy(io.Discard, response.Body)
			_ = response.Body.Close()
			statuses <- response.StatusCode
		}(a)
	}
	wg.Wait()
	x, y := <-statuses, <-statuses
	if !((x == 201 && y == 409) || (x == 409 && y == 201)) {
		t.Fatalf("конкурирующие занятия: %d, %d", x, y)
	}
	apiRequest(t, router, "POST", "/api/auth/token/logout", token, nil, 204)
	apiRequest(t, router, "GET", "/api/groups", token, nil, 401)
	if err := json.Unmarshal(apiRequest(t, router, "POST", "/api/auth/token/login", "", credentials, 200), &login); err != nil {
		t.Fatal(err)
	}
	if err := db.Exec("UPDATE api_tokens SET expires_at = NOW() - INTERVAL '1 day'").Error; err != nil {
		t.Fatal(err)
	}
	apiRequest(t, router, "GET", "/api/auth/users/me", login.Token, nil, 401)
}

func TestCRUDResources(t *testing.T) {
	db := testDatabase(t)
	router := testRouter(db)
	credentials := map[string]any{"username": "crud_test", "password": "password123"}
	apiRequest(t, router, "POST", "/api/auth/users", "", credentials, 201)
	var login struct {
		Token string `json:"auth_token"`
	}
	if err := json.Unmarshal(apiRequest(t, router, "POST", "/api/auth/token/login", "", credentials, 200), &login); err != nil {
		t.Fatal(err)
	}
	token := login.Token
	ids := make(map[string]int64)
	check := func(resource string, body, update map[string]any, changedField string) {
		t.Helper()
		var created struct {
			ID int64 `json:"id"`
		}
		if err := json.Unmarshal(apiRequest(t, router, "POST", "/api/"+resource, token, body, 201), &created); err != nil {
			t.Fatal(err)
		}
		ids[resource] = created.ID
		url := fmt.Sprintf("/api/%s/%d", resource, created.ID)
		apiRequest(t, router, "GET", "/api/"+resource, token, nil, 200)
		apiRequest(t, router, "GET", url, token, nil, 200)
		apiRequest(t, router, "PUT", url, token, update, 204)
		var fetched map[string]any
		if err := json.Unmarshal(apiRequest(t, router, "GET", url, token, nil, 200), &fetched); err != nil {
			t.Fatal(err)
		}
		// JSON возвращает числа как float64, поэтому сравниваем сериализованное поле.
		want, _ := json.Marshal(update[changedField])
		got, _ := json.Marshal(fetched[changedField])
		if !bytes.Equal(want, got) {
			t.Fatalf("%s: не обновлено поле %s", resource, changedField)
		}
	}
	check("groups", map[string]any{"name": "К1001", "course": 1}, map[string]any{"name": "К2001", "course": 2}, "course")
	check("classrooms", map[string]any{"number": "101", "building": "А"}, map[string]any{"number": "102", "building": "А"}, "number")
	check("subjects", map[string]any{"name": "Алгебра"}, map[string]any{"name": "Математика"}, "name")
	check("semesters", map[string]any{"academic_year": "2025-2026", "number": 1}, map[string]any{"academic_year": "2026-2027", "number": 1}, "academic_year")
	check("teachers", map[string]any{"first_name": "Иван", "last_name": "Петров", "fixed_classroom_id": ids["classrooms"]}, map[string]any{"first_name": "Иван", "last_name": "Сидоров", "fixed_classroom_id": nil, "is_active": true}, "fixed_classroom_id")
	check("students", map[string]any{"first_name": "Олег", "last_name": "Петров", "group_id": ids["groups"], "birth_date": "2004-01-01"}, map[string]any{"first_name": "Олег", "last_name": "Сидоров", "group_id": ids["groups"], "birth_date": nil, "is_active": true}, "birth_date")
	teachingBody := map[string]any{"teacher_id": ids["teachers"], "subject_id": ids["subjects"], "group_id": ids["groups"], "semester_id": ids["semesters"]}
	check("teachings", teachingBody, teachingBody, "teacher_id")
	check("grades", map[string]any{"student_id": ids["students"], "subject_id": ids["subjects"], "semester_id": ids["semesters"], "value": 5}, map[string]any{"student_id": ids["students"], "subject_id": ids["subjects"], "semester_id": ids["semesters"], "value": 2}, "value")
	check("schedules", map[string]any{"assignment_id": ids["teachings"], "classroom_id": ids["classrooms"], "weekday": 1, "lesson_number": 1}, map[string]any{"assignment_id": ids["teachings"], "classroom_id": ids["classrooms"], "weekday": 2, "lesson_number": 1}, "weekday")
	apiRequest(t, router, "GET", "/api/groups/0", token, nil, 400)
	apiRequest(t, router, "POST", "/api/groups", token, map[string]any{"name": "Invalid", "course": 0}, 400)
	apiRequest(t, router, "POST", "/api/semesters", token, map[string]any{"academic_year": "2026-2028", "number": 1}, 400)
	apiRequest(t, router, "POST", "/api/subjects", token, map[string]any{"name": "Subject", "unexpected": true}, 400)
	apiRequest(t, router, "DELETE", fmt.Sprintf("/api/groups/%d", ids["groups"]), token, nil, 409)
	for _, resource := range []string{"schedules", "grades", "teachings", "students", "teachers", "semesters", "subjects", "classrooms", "groups"} {
		url := fmt.Sprintf("/api/%s/%d", resource, ids[resource])
		apiRequest(t, router, "DELETE", url, token, nil, 204)
		apiRequest(t, router, "GET", url, token, nil, 404)
		apiRequest(t, router, "DELETE", url, token, nil, 404)
	}
}
