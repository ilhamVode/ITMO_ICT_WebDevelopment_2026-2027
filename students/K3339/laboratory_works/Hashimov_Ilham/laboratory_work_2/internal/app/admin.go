package app

import (
	"html/template"
	"net/http"
	"strconv"
	"time"
)

type AdminView struct {
	userRepo       UserRepository
	homeworkRepo   HomeworkRepository
	submissionRepo SubmissionRepository
}

func NewAdminView(
	userRepo UserRepository,
	homeworkRepo HomeworkRepository,
	submissionRepo SubmissionRepository,
) *AdminView {
	return &AdminView{
		userRepo:       userRepo,
		homeworkRepo:   homeworkRepo,
		submissionRepo: submissionRepo,
	}
}

func (v *AdminView) Index(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("web/templates/admin/index.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	err = tmpl.Execute(w, nil)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		return
	}
}

func (v *AdminView) Users(w http.ResponseWriter, r *http.Request) {
	users := []User{}

	users, err := v.userRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/admin/users.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, users)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		return
	}
}
func (v *AdminView) Homeworks(w http.ResponseWriter, r *http.Request) {
	homeworks := []Homework{}

	homeworks, err := v.homeworkRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/admin/homeworks.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, homeworks)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		return
	}
}

func (v *AdminView) Submissions(w http.ResponseWriter, r *http.Request) {
	submissions := []Submission{}

	submissions, err := v.submissionRepo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/admin/submissions.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}
	err = tmpl.Execute(w, submissions)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		return
	}
}

func (v *AdminView) CreateUser(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/admin/user_create.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()

		if err != nil {
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			return
		}

		username := r.FormValue("username")
		password := r.FormValue("password")
		className := r.FormValue("class_name")
		passportNumber := r.FormValue("passport_number")
		homeAddress := r.FormValue("home_address")
		nationality := r.FormValue("nationality")

		isAdmin := r.FormValue("is_admin") == "on"

		user := User{
			Username:       username,
			PasswordHash:   password,
			IsAdmin:        isAdmin,
			ClassName:      className,
			PassportNumber: passportNumber,
			HomeAddress:    homeAddress,
			Nationality:    nationality,
		}

		_, err = v.userRepo.Create(r.Context(), user)
		if err != nil {
			http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/users", http.StatusSeeOther)
		return
	}
}
func (v *AdminView) CreateHomework(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/admin/homework_create.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}
		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
			return
		}
		return
	}
	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			return
		}

		subject := r.FormValue("subject")
		teacherIDStr := r.FormValue("teacher_id")
		startDateStr := r.FormValue("start_date")
		endDateStr := r.FormValue("end_date")
		text := r.FormValue("text")
		penalty := r.FormValue("penalty")

		teacherID, err := strconv.ParseInt(teacherIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Некорректный ID преподавателя", http.StatusBadRequest)
			return
		}

		startDate, err := time.Parse("2006-01-02T15:04", startDateStr)
		if err != nil {
			http.Error(w, "Некорректная дата начала", http.StatusBadRequest)
			return
		}

		endDate, err := time.Parse("2006-01-02T15:04", endDateStr)
		if err != nil {
			http.Error(w, "Некорректная дата окончания", http.StatusBadRequest)
			return
		}

		homework := Homework{
			Subject:   subject,
			TeacherID: teacherID,
			IssueDate: time.Now(),
			StartDate: startDate,
			EndDate:   endDate,
			Text:      text,
			Penalty:   penalty,
		}

		_, err = v.homeworkRepo.Create(r.Context(), homework)
		if err != nil {
			http.Error(w, "Ошибка создания домашнего задания", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/homeworks", http.StatusSeeOther)
		return
	}
}
func (v *AdminView) CreateSubmission(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/admin/submission_create.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, nil)
		if err != nil {
			http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
			return
		}

		return
	}

	if r.Method == http.MethodPost {
		err := r.ParseForm()
		if err != nil {
			http.Error(w, "Ошибка при парсинге формы", http.StatusBadRequest)
			return
		}

		studentIDStr := r.FormValue("student_id")
		homeworkIDStr := r.FormValue("homework_id")
		text := r.FormValue("text")

		studentID, err := strconv.ParseInt(studentIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Некорректный ID студента", http.StatusBadRequest)
			return
		}

		homeworkID, err := strconv.ParseInt(homeworkIDStr, 10, 64)
		if err != nil {
			http.Error(w, "Некорректный ID домашнего задания", http.StatusBadRequest)
			return
		}

		submission := Submission{
			StudentID:   studentID,
			HomeworkID:  homeworkID,
			Text:        text,
			SubmittedAt: time.Now(),
			Grade:       nil,
		}

		_, err = v.submissionRepo.Create(r.Context(), submission)
		if err != nil {
			http.Error(w, "Ошибка создания сдачи", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/admin/submissions", http.StatusSeeOther)
		return
	}
}
