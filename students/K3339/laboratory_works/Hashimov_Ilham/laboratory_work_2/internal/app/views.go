package app

import (
	"errors"
	"html/template"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
)

type HomeworkView struct {
	repo HomeworkRepository
	auth *Auth
}

type HomeworkPageData struct {
	Homework Homework
	User     User
	IsLogged bool
}

type HomeworksPageData struct {
	Homeworks  []Homework
	User       User
	IsLogged   bool
	Page       int
	TotalPages int
	HasPrev    bool
	HasNext    bool
	PrevPage   int
	NextPage   int
	Search     string
	Subject    string
}

type UserView struct {
	repo UserRepository
}

func NewHomeworkView(repo HomeworkRepository, auth *Auth) *HomeworkView {
	return &HomeworkView{
		repo: repo,
		auth: auth,
	}
}
func NewUserView(repo UserRepository) *UserView {
	return &UserView{repo: repo}
}

func (v *HomeworkView) Detail(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	homework, err := v.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Домашнее задание не найдено", http.StatusNotFound)
			return
		}

		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/homework.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	user, isLogged := v.auth.CurrentUser(r)

	data := HomeworkPageData{
		Homework: homework,
		User:     user,
		IsLogged: isLogged,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		return
	}
}

func (v *HomeworkView) List(w http.ResponseWriter, r *http.Request) {
	const limit = 5

	page := 1

	pageStr := r.URL.Query().Get("page")
	if pageStr != "" {
		parsedPage, err := strconv.Atoi(pageStr)
		if err == nil && parsedPage > 0 {
			page = parsedPage
		}
	}

	search := r.URL.Query().Get("search")
	subject := r.URL.Query().Get("subject")

	count, err := v.repo.CountFiltered(r.Context(), search, subject)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	totalPages := int((count + int64(limit) - 1) / int64(limit))

	if totalPages == 0 {
		totalPages = 1
	}

	if page > totalPages {
		page = totalPages
	}

	offset := (page - 1) * limit

	homeworks, err := v.repo.GetPageFiltered(
		r.Context(),
		limit,
		offset,
		search,
		subject,
	)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/homeworks.html")
	if err != nil {
		http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
		return
	}

	user, isLogged := v.auth.CurrentUser(r)

	data := HomeworksPageData{
		Homeworks:  homeworks,
		User:       user,
		IsLogged:   isLogged,
		Page:       page,
		TotalPages: totalPages,
		HasPrev:    page > 1,
		HasNext:    page < totalPages,
		PrevPage:   page - 1,
		NextPage:   page + 1,
		Search:     search,
		Subject:    subject,
	}

	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
		return
	}
}

func (v *HomeworkView) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	homework, err := v.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Домашнее задание не найдено", http.StatusNotFound)
			return
		}

		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/homework_update.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, homework)
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

		homework.Subject = subject
		homework.TeacherID = teacherID
		homework.StartDate = startDate
		homework.EndDate = endDate
		homework.Text = text
		homework.Penalty = penalty

		err = v.repo.Update(r.Context(), homework)
		if err != nil {
			http.Error(w, "Ошибка обновления домашнего задания", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/homework/"+idStr, http.StatusSeeOther)
		return
	}
}

func (v *HomeworkView) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.PathValue("id")

	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Некорректный ID", http.StatusBadRequest)
		return
	}

	homework, err := v.repo.GetByID(r.Context(), id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Домашнее задание не найдено", http.StatusNotFound)
			return
		}

		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/homework_delete.html")
		if err != nil {
			http.Error(w, "Ошибка загрузки шаблона", http.StatusInternalServerError)
			return
		}

		err = tmpl.Execute(w, homework)
		if err != nil {
			http.Error(w, "Ошибка отображения шаблона", http.StatusInternalServerError)
			return
		}

		return
	}

	if r.Method == http.MethodPost {
		err := v.repo.Delete(r.Context(), id)
		if err != nil {
			http.Error(w, "Ошибка удаления домашнего задания", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/homeworks", http.StatusSeeOther)
		return
	}
}

func (v *HomeworkView) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/homework_create.html")
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

		if endDate.Before(startDate) {
			http.Error(w, "Дата окончания не может быть раньше даты начала", http.StatusBadRequest)
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

		_, err = v.repo.Create(r.Context(), homework)
		if err != nil {
			http.Error(w, "Ошибка создания домашнего задания", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/homeworks", http.StatusSeeOther)
		return
	}
}

func (v *UserView) List(w http.ResponseWriter, r *http.Request) {
	users, err := v.repo.GetAll(r.Context())
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

	tmpl, err := template.ParseFiles("web/templates/users.html")
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

func (v *UserView) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		tmpl, err := template.ParseFiles("web/templates/user_create.html")
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

		user := User{
			Username:       username,
			PasswordHash:   password,
			IsAdmin:        false,
			ClassName:      className,
			PassportNumber: passportNumber,
			HomeAddress:    homeAddress,
			Nationality:    nationality,
		}

		_, err = v.repo.Create(r.Context(), user)
		if err != nil {
			http.Error(w, "Ошибка создания пользователя", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/users", http.StatusSeeOther)
		return
	}
}
