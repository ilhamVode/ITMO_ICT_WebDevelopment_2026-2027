# Лабораторная работа №3

## Вариант 12. Учебная часть колледжа

В третьей работе реализован REST API для учебной части колледжа. В нём хранятся группы, студенты, преподаватели, дисциплины, кабинеты, назначения, расписание и семестровые оценки. Работа на Go, с GORM и PostgreSQL.

Ниже приведены запросы и ответы на тестовых данных в PostgreSQL. Значение токена заменено на `<token>`. В другой базе ID записей могут отличаться.

---

## Task 1. Проект

Проект находится в `laboratory_work_3`. Основные файлы:

```text
laboratory_work_3/
├── cmd/server/
│   ├── main.go
│   ├── urls.go
│   └── api_test.go
├── internal/
│   ├── student/, group/, teacher/, classroom/, subject/
│   ├── semester/, teaching/, schedule/, grade/
│   ├── analytics/, report/, auth/
│   ├── httpapi/http.go
│   └── database/
├── migrations/
│   ├── init.sql
│   └── business_rules.sql
├── .env.example
├── .gitignore
├── go.mod
└── go.sum
```

В папке каждой сущности находятся `model.go`, `repository.go`, `service.go` и `handler.go`. Модель описывает данные, репозиторий выполняет запросы к БД, сервис проверяет данные, обработчик принимает HTTP-запрос и формирует ответ.

В `main.go` подключение создаётся через `database.NewPostgres(dsn)`. Затем `database.InitSchema(db)` применяет оба SQL-файла в одной транзакции. Таблицы создаются через `CREATE TABLE IF NOT EXISTS`, поэтому повторный запуск не удаляет записи. Функции и триггеры обновляются через `CREATE OR REPLACE`.

Для работы с моделями и запросами используется GORM.

Проверка доступности сервера:

```powershell
curl.exe http://localhost:8080/health
```

```text
GET /health
HTTP 200 OK
{
  "status": "ok"
}
```

---

## Task 2. Модель данных

| Таблица | Основные поля | Назначение |
| --- | --- | --- |
| `study_groups` | `name`, `course` | Учебная группа и курс |
| `students` | Имя, фамилия, `birth_date`, `group_id`, `is_active` | Студент и его группа |
| `classrooms` | `number`, `building` | Кабинет и корпус |
| `teachers` | ФИО, `fixed_classroom_id`, `is_active` | Преподаватель; собственного кабинета может не быть |
| `subjects` | `name` | Дисциплина |
| `semesters` | `academic_year`, `number` | Учебный год и семестр |
| `teaching_assignments` | `teacher_id`, `subject_id`, `group_id`, `semester_id` | Кто, что и какой группе преподаёт за семестр |
| `schedules` | `assignment_id`, `classroom_id`, `weekday`, `lesson_number` | Занятие в расписании |
| `grades` | `student_id`, `subject_id`, `semester_id`, `value` | Итоговая оценка за семестр |
| `api_users` | `username`, `password_hash` | Пользователь API |
| `api_tokens` | `key`, `user_id`, `expires_at` | Хеш токена, пользователь и срок действия |

У каждой предметной сущности есть первичный ключ `id`. Связи заданы внешними ключами в `init.sql`:

```text
study_groups ──< students
teachers ──< teaching_assignments >── subjects
study_groups ──< teaching_assignments >── semesters
teaching_assignments ──< schedules >── classrooms
students ──< grades >── subjects
semesters ──< grades
classrooms ──< teachers.fixed_classroom_id
api_users ──< api_tokens
```

Здесь `──<` обозначает связь «один ко многим». Например, у одной группы может быть несколько студентов, а каждый студент относится к одной группе. Основные связи модели:

| Связь | Тип | Как хранится в БД |
| --- | --- | --- |
| Группа → студенты | 1:m | `students.group_id` ссылается на `study_groups.id` |
| Преподаватель, дисциплина, группа, семестр → назначения | Каждая 1:m | Четыре внешних ключа в `teaching_assignments` |
| Назначение → занятия; кабинет → занятия | Каждая 1:m | `schedules.assignment_id` и `schedules.classroom_id` |
| Студент, дисциплина, семестр → оценки | Каждая 1:m | Три внешних ключа в `grades` |
| Кабинет → закреплённые преподаватели | 1:m, кабинет необязателен | `teachers.fixed_classroom_id` допускает `NULL` |
| Пользователь → токены | 1:m | `api_tokens.user_id` |
| Преподаватели ↔ дисциплины; преподаватели ↔ группы; группы ↔ дисциплины | Каждая m:m | Через назначения `teaching_assignments`, с учётом семестра |
| Студенты ↔ дисциплины, по которым выставлены оценки | m:m | Через `grades`, где также хранятся семестр и значение оценки |

Преподаватели и дисциплины связаны m:m: преподаватель ведёт несколько дисциплин, одну дисциплину ведут разные преподаватели. Промежуточная модель `TeachingAssignment` также хранит группу и семестр.

Пример модели:

```go
type TeachingAssignment struct {
    ID         int64 `gorm:"primaryKey;column:id" json:"id"`
    TeacherID  int64 `gorm:"column:teacher_id" json:"teacher_id"`
    SubjectID  int64 `gorm:"column:subject_id" json:"subject_id"`
    GroupID    int64 `gorm:"column:group_id" json:"group_id"`
    SemesterID int64 `gorm:"column:semester_id" json:"semester_id"`

    Subject *subject.Subject `gorm:"foreignKey:SubjectID;references:ID" json:"subject,omitempty"`
}
```

Теги `gorm` задают связь с таблицей и колонками. Теги `json` задают названия полей ответа. Указатель, например `*int64` у `FixedClassroomID`, позволяет хранить `nil`, то есть SQL `NULL`: у преподавателя может не быть кабинета.

В БД нельзя повторно создать группу с тем же названием, кабинет с той же парой «номер + корпус», одинаковое назначение или вторую оценку студенту по той же дисциплине за тот же семестр.

---

## Практика 3.1. Запросы ORM и связи

Запросы к моделям выполняются через GORM:

```go
err := db.WithContext(ctx).Create(&student).Error
err := db.WithContext(ctx).First(&student, id).Error
err := db.WithContext(ctx).Order("id ASC").Find(&students).Error
err := db.WithContext(ctx).Where("group_id = ?", groupID).Find(&students).Error
```

`Create` добавляет запись, `First` получает студента по ID, `Find` получает список, `Where` выбирает студентов указанной группы.

Список студентов группы задан полем `Students`:

```go
Students []student.Student `gorm:"foreignKey:GroupID;references:ID" json:"students,omitempty"`
```

Связанные студенты загружаются через `Preload("Students")`, дисциплины преподавателя — через `Preload("Assignments.Subject")`.

Агрегационный запрос считает студентов по курсам:

```sql
SELECT g.course, COUNT(DISTINCT st.id) AS student_count
FROM schedules sc
JOIN teaching_assignments ta ON ta.id = sc.assignment_id
JOIN study_groups g ON g.id = ta.group_id
JOIN students st ON st.group_id = g.id
WHERE sc.classroom_id = ? AND ta.semester_id = ? AND st.is_active = TRUE
GROUP BY g.course
ORDER BY g.course;
```

Он находится в `internal/analytics/repository.go`. `DISTINCT` исключает повторный подсчёт студента, если у его группы несколько занятий в одном кабинете. Значения передаются параметрами, на места `?`.

---

## Практика 3.2. REST API и сериализация

Общий путь запроса:

```text
HTTP-клиент → ServeMux → проверка токена → Handler → Service → Repository → PostgreSQL
                              ← JSON-ответ ←
```

Маршруты связываются с методами обработчиков:

```go
mux.HandleFunc("GET /api/students/{id}", h.GetByID)
mux.HandleFunc("POST /api/students", h.Create)
mux.HandleFunc("PUT /api/students/{id}", h.Update)
mux.HandleFunc("DELETE /api/students/{id}", h.Delete)
```

Обработчик читает JSON в структуру запроса, сервис проверяет данные и сохраняет их через репозиторий. Ответ формируется с помощью `encoding/json`. Пример создания оценки из `internal/grade/handler.go`:

```go
func (h *GradeHandler) Create(w http.ResponseWriter, r *http.Request) {
    var req gradeRequest

    if err := httpapi.ReadJSON(w, r, &req); err != nil {
        httpapi.WriteError(w, http.StatusBadRequest, err.Error())
        return
    }

    grade := req.ToModel()
    if err := h.service.Create(r.Context(), &grade); err != nil {
        httpapi.HandleError(w, err, ErrInvalidGrade)
        return
    }

    w.Header().Set("Location", "/api/grades/"+strconv.FormatInt(grade.ID, 10))
    httpapi.WriteJSON(w, http.StatusCreated, grade)
}
```

Неизвестные поля и лишний JSON-объект в теле запроса отклоняются. Даты рождения принимаются в формате `YYYY-MM-DD`.

### CRUD-эндпоинты

Для каждого ресурса реализованы пять операций: `GET /api/<ресурс>` — список, `GET /api/<ресурс>/{id}` — объект, `POST` на список — создание, `PUT` по ID — полное обновление, `DELETE` по ID — удаление.

| Ресурс | Поля JSON при создании и обновлении |
| --- | --- |
| `groups` | `name`, `course` |
| `students` | `first_name`, `last_name`, `group_id`, необязательный `birth_date` |
| `teachers` | `first_name`, `last_name`, необязательные `middle_name`, `fixed_classroom_id` |
| `classrooms` | `number`, `building` |
| `subjects` | `name` |
| `semesters` | `academic_year` в виде `2026-2027`, `number` — 1 или 2 |
| `teachings` | `teacher_id`, `subject_id`, `group_id`, `semester_id` |
| `schedules` | `assignment_id`, `classroom_id`, `weekday` — 1–7, `lesson_number` |
| `grades` | `student_id`, `subject_id`, `semester_id`, `value` — 2–5 |

У студентов и преподавателей при создании `is_active` устанавливается в `true`. При `PUT` это поле обязательно передаётся явно. `null` можно передать для необязательной даты рождения, отчества и закреплённого кабинета. Частичного обновления через `PATCH` в текущем API нет.

Чтение возвращает `200`, создание — `201` и заголовок `Location`, обновление и удаление — `204` без тела. Некорректный запрос получает `400`, отсутствующий объект — `404`, повторная запись или конфликт расписания — `409`.

### Создание объектов

Для примеров подготовлены группы К3339, К2339, К1339; два преподавателя; дисциплины «Математика» и «Физика»; кабинеты 301 и 302; два семестра. В группе К3339 находятся Ильхам Хашимов и Анна Иванова. Все эти объекты создаются через API.

```text
POST /api/groups
{"course":3,"name":"К3339"}
HTTP 201 Created
{"id":1,"name":"К3339","course":3}
```

```text
POST /api/students
{"birth_date":"2005-01-01","first_name":"Ильхам","group_id":1,"last_name":"Хашимов"}
HTTP 201 Created
{"birth_date":"2005-01-01","first_name":"Ильхам","group_id":1,"id":1,"is_active":true,"last_name":"Хашимов"}
```

Назначение связывает преподавателя, дисциплину, группу и семестр. Расписание ссылается уже на это назначение:

```text
POST /api/teachings
{"group_id":1,"semester_id":1,"subject_id":1,"teacher_id":1}
HTTP 201 Created
{"id":1,"teacher_id":1,"subject_id":1,"group_id":1,"semester_id":1}
```

```text
POST /api/schedules
{"assignment_id":1,"classroom_id":1,"lesson_number":1,"weekday":1}
HTTP 201 Created
{"id":1,"assignment_id":1,"classroom_id":1,"weekday":1,"lesson_number":1}
```

---

## Авторизация по токенам

Регистрация, вход, текущий пользователь и выход реализованы в `internal/auth`. При выходе переданный токен удаляется.

| Метод и адрес | Что делает |
| --- | --- |
| `POST /api/auth/users` | Регистрация; тело содержит `username`, `password` |
| `POST /api/auth/token/login` | Вход с теми же полями, возвращает `auth_token` |
| `GET /api/auth/users/me` | Данные текущего пользователя по токену |
| `POST /api/auth/token/logout` | Удаление переданного токена |

Имя пользователя содержит 3–50 байт, буквы, цифры, `_` или `-`; пароль — 8–72 байта. Пароль сохраняется как bcrypt-хеш. Токен создаётся из 32 случайных байтов, действует 24 часа, в БД хранится только его SHA-256-хеш.

```text
POST /api/auth/users
{"password":"password123","username":"student_test"}
HTTP 201 Created
{"id":1,"username":"student_test"}
```

```text
POST /api/auth/token/login
{"password":"password123","username":"student_test"}
HTTP 200 OK
{"auth_token":"<token>"}
```

В PowerShell токен можно получить в переменную и затем передавать заголовком:

```powershell
$base = "http://localhost:8080"
$credentials = '{"username":"student_test","password":"password123"}'
$login = Invoke-RestMethod -Method Post -Uri "$base/api/auth/token/login" -ContentType "application/json" -Body $credentials
$headers = @{ Authorization = "Token $($login.auth_token)" }
Invoke-RestMethod -Uri "$base/api/auth/users/me" -Headers $headers
```

```text
GET /api/auth/users/me
HTTP 200 OK
{"id":1,"username":"student_test"}
```

Проверка токена подключена к учебным эндпоинтам через `RequireAuth`. Каждый запрос передаёт `Authorization: Token <token>`; также поддерживается `Bearer`. Зарегистрированные пользователи имеют одинаковый доступ к учебным данным. `/health`, регистрация и вход публичные.

После выхода запрос с прежним токеном отклоняется:

```text
POST /api/auth/token/logout
HTTP 204 No Content
```

```text
GET /api/groups
HTTP 401 Unauthorized
{"error":"требуется действующий токен или верные данные для входа"}
```

---

## Вложенные объекты

### Один ко многим: группа и студенты

Связь реализована в `migrations/init.sql`: внешний ключ `students.group_id` указывает на `study_groups.id`. В `internal/student/model.go` это поле `GroupID`, а обратная сторона связи — список `Students` в `internal/group/model.go`. В GORM тег `foreignKey:GroupID;references:ID` указывает, по каким полям связаны записи.

При создании студента передаётся `group_id`, поэтому он появляется в списке студентов этой группы. Несуществующий ID группы не проходит проверку внешнего ключа. Чтобы перевести студента, через `PUT /api/students/{id}` меняется `group_id` вместе с остальными полями.

Для `GET /api/groups/{id}` репозиторий использует:

```go
err := r.db.WithContext(ctx).
    Preload("Students").
    First(&group, id).Error
```

```text
GET /api/groups/1
HTTP 200 OK
{
  "id": 1,
  "name": "К3339",
  "course": 3,
  "students": [
    {
      "id": 1,
      "first_name": "Ильхам",
      "last_name": "Хашимов",
      "birth_date": "2005-01-01",
      "group_id": 1,
      "is_active": true
    },
    {
      "id": 2,
      "first_name": "Анна",
      "last_name": "Иванова",
      "birth_date": null,
      "group_id": 1,
      "is_active": true
    }
  ]
}
```

Один HTTP-запрос возвращает группу вместе со студентами.

### Многие ко многим: преподаватель и дисциплины

В `internal/teaching/model.go` связующая модель содержит `TeacherID` и `SubjectID`. В БД это внешние ключи на `teachers.id` и `subjects.id`. Каждая строка `teaching_assignments` связывает одного преподавателя с одной дисциплиной для конкретной группы и семестра. Несколько строк дают связь m:m.

Назначения из тестового примера:

| `teacher_id` | `subject_id` | `group_id` | `semester_id` |
| --- | --- | --- | --- |
| 1 | 1 | 1 | 1 |
| 2 | 1 | 2 | 1 |
| 1 | 2 | 1 | 1 |
| 1 | 1 | 3 | 1 |
| 1 | 1 | 1 | 2 |

Преподаватель 1, Иван Петров, ведёт математику и физику. Математику ведут преподаватели 1 и 2, то есть Иван Петров и Мария Соколова. Это показывает обе стороны m:m. Ограничение `UNIQUE (teacher_id, subject_id, group_id, semester_id)` запрещает только полное повторение назначения: ту же дисциплину можно вести в другой группе или семестре.

Связь добавляется через `POST /api/teachings`, изменяется через `PUT /api/teachings/{id}` и удаляется через `DELETE /api/teachings/{id}`. При удалении назначения преподаватель и дисциплина остаются. Если на назначение ссылается расписание, сначала нужно удалить его занятия.

`GET /api/teachers/{id}` возвращает назначения с вложенной дисциплиной и отдельный список `subjects`. Используется `Preload("Assignments.Subject")`; список дисциплин собирается без повторов по `subject_id`.

В `internal/teacher/model.go` и `internal/teacher/repository.go` загрузка устроена так:

```go
// Поле модели Teacher — обратная связь с назначениями.
Assignments []teaching.TeachingAssignment `gorm:"foreignKey:TeacherID;references:ID" json:"assignments,omitempty"`

// Запрос репозитория — назначения и дисциплина каждого назначения.
err := r.db.WithContext(ctx).
    Preload("Assignments.Subject").
    First(&teacher, id).Error
```

Из ответа выведены поля `id`, `first_name`, `last_name` и `subjects`:

```powershell
Invoke-RestMethod -Uri "$base/api/teachers/1" -Headers $headers |
    Select-Object id, first_name, last_name, subjects | ConvertTo-Json -Depth 6
```

```text
GET /api/teachers/1
HTTP 200 OK
{
  "id": 1,
  "first_name": "Иван",
  "last_name": "Петров",
  "subjects": [
    {
      "id": 1,
      "name": "Математика"
    },
    {
      "id": 2,
      "name": "Физика"
    }
  ]
}
```

Математика в назначениях этого преподавателя встречается в нескольких группах и семестрах, а в `subjects` выводится один раз. Тег `gorm:"-"` исключает поле `subjects` из сохраняемых полей модели.

Также `GET /api/schedules/{id}` возвращает назначение с дисциплиной и кабинет, а `GET /api/grades/{id}` — студента, дисциплину и семестр.

---

## Task 3. Аналитические запросы варианта 12

Для запросов указывается семестр. День недели: 1 — понедельник, 7 — воскресенье. Отсутствующие или некорректные параметры возвращают `400`; если подходящих данных нет, возвращается пустой массив.

### 1. Предмет в группе на указанной паре

```text
GET /api/analytics/lesson?group_id=1&semester_id=1&weekday=1&lesson_number=1
HTTP 200 OK
[
  {
    "schedule_id": 1,
    "lesson_number": 1,
    "subject_id": 1,
    "subject_name": "Математика",
    "teacher_id": 1,
    "teacher_first_name": "Иван",
    "teacher_last_name": "Петров",
    "classroom_id": 1,
    "classroom_number": "301",
    "classroom_building": "Главный"
  }
]
```

### 2. Преподаватели заданной группы

```text
GET /api/analytics/teachers?group_id=1&semester_id=1
HTTP 200 OK
[
  {
    "id": 1,
    "first_name": "Иван",
    "last_name": "Петров"
  }
]
```

В запросе используется `SELECT DISTINCT`, поэтому Иван Петров выводится один раз, хотя ведёт у группы две дисциплины.

### 3. Группы, где преподаватель ведёт заданный предмет

```text
GET /api/analytics/groups?teacher_id=1&subject_id=1&semester_id=1
HTTP 200 OK
[
  {
    "id": 1,
    "name": "К3339",
    "course": 3
  },
  {
    "id": 3,
    "name": "К1339",
    "course": 1
  }
]
```

### 4. Расписание группы на день недели

```text
GET /api/analytics/schedule?group_id=1&semester_id=1&weekday=1
HTTP 200 OK
```

Ответ содержит занятия с теми же полями, что и в первом запросе. Результат в виде таблицы:

| ID занятия | Пара | Дисциплина | Преподаватель | Кабинет |
| --- | --- | --- | --- | --- |
| 1 | 1 | Математика | Иван Петров | 301 |
| 7 | 2 | Физика | Иван Петров | 301 |

Строки упорядочены по номеру пары. Пропуски в ID возможны: PostgreSQL расходует значения последовательности и при попытках вставки, завершившихся ошибкой.

### 5. Количество студентов по курсам в указанном кабинете

В условии используется выражение «в указанном классе». В этой реализации оно трактуется как учебный кабинет: считаются активные студенты групп, у которых запланированы занятия в этом кабинете за выбранный семестр. Фактическая посещаемость не хранится.

```text
GET /api/analytics/student-counts?classroom_id=1&semester_id=1
HTTP 200 OK
[
  {
    "course": 3,
    "student_count": 2
  }
]
```

У К3339 две пары в кабинете 301, но два студента посчитаны один раз благодаря `COUNT(DISTINCT st.id)`.

---

## Task 4. Семестровые оценки и ведомость

Оценка создаётся отдельным запросом:

```text
POST /api/grades
{"semester_id":1,"student_id":1,"subject_id":1,"value":5}
HTTP 201 Created
{"id":1,"student_id":1,"subject_id":1,"semester_id":1,"value":5}
```

Ведомость получается через `GET /api/reports/statement?group_id=1&semester_id=1`. Полный JSON содержит группу, семестр, студентов с оценками, средние баллы по студентам и дисциплинам, количество студентов и оценок, средний балл группы.

На тестовых данных результат такой:

| Студент | Математика | Физика | Средний балл |
| --- | --- | --- | --- |
| Анна Иванова | 3 | нет оценки | 3 |
| Ильхам Хашимов | 5 | 4 | 4,5 |

Средний по математике — 4, по физике — 4. Средний по группе считается по всем выставленным оценкам: `(3 + 5 + 4) / 3 = 4`.

Из ответа можно вывести только итоговые поля:

```powershell
Invoke-RestMethod -Uri "$base/api/reports/statement?group_id=1&semester_id=1" -Headers $headers |
    Select-Object student_count, grade_count, average | ConvertTo-Json
```

```json
{
  "student_count": 2,
  "grade_count": 3,
  "average": 4
}
```

Отсутствующая оценка возвращается как `null` и не участвует в вычислениях. Если оценок нет совсем, средний балл тоже `null`. Значения округляются до двух знаков. Для чтения частей ведомости используется одна транзакция `Repeatable Read`.

### Выгрузка документа

Параметр `format=csv` возвращает файл ведомости:

```powershell
curl.exe "$base/api/reports/statement?group_id=1&semester_id=1&format=csv" `
    -H "Authorization: Token $($login.auth_token)" -o statement.csv
```

Полученное содержимое CSV:

```text
Группа;К3339
Учебный год;2026-2027;Семестр;1
ID;Фамилия;Имя;Обучается;Математика;Физика;Средний балл
2;Иванова;Анна;true;3;;3.00
1;Хашимов;Ильхам;true;5;4;4.50
;Средний балл по дисциплине;;;4.00;4.00;4.00
Количество студентов;2
Количество оценок;3
```

В ответе устанавливаются `Content-Type: text/csv; charset=utf-8` и `Content-Disposition: attachment; filename=statement_1_1.csv`. Разделитель — `;`, добавлен UTF-8 BOM для русского текста в Excel.

### Изменение оценки

Меняем оценку по математике с 5 на 2:

```text
PUT /api/grades/1
{"semester_id":1,"student_id":1,"subject_id":1,"value":2}
HTTP 204 No Content
```

Повторный запрос ведомости возвращает средний балл группы 3: `(3 + 2 + 4) / 3 = 3`. Данные за другой семестр в расчёт не попадают.

---

## Проверки данных, отчисление и увольнение

В сервисах проверяются обязательные поля, положительные ID, курс и номер пары, формат учебного года, диапазоны дня недели и оценки. Дополнительные проверки находятся в `migrations/business_rules.sql`.

В одном семестре на одной паре не могут одновременно быть заняты одна группа, один преподаватель или один кабинет. Проверяются создание и изменение расписания, а также изменение назначения, которое уже включено в расписание.

Попытка поставить физику группе К3339 на занятую первую пару:

```text
POST /api/schedules
{"assignment_id":3,"classroom_id":2,"lesson_number":1,"weekday":1}
HTTP 409 Conflict
{"error":"в эту пару заняты группа, преподаватель или кабинет"}
```

Проверка выполняется триггером PostgreSQL. `pg_advisory_xact_lock` последовательно пропускает конкурирующие изменения расписания и назначений, чтобы два одновременных запроса не заняли одну пару.

Оценку можно выставить только активному студенту по дисциплине, назначенной его группе за нужный семестр:

```text
POST /api/grades
{"semester_id":1,"student_id":3,"subject_id":2,"value":4}
HTTP 400 Bad Request
{"error":"дисциплина не назначена группе студента в этом семестре"}
```

```text
POST /api/students/2/expel
HTTP 204 No Content
```

```text
POST /api/teachers/2/dismiss
HTTP 204 No Content
```

`POST /api/students/{id}/expel` устанавливает `is_active=false` и сохраняет студента с оценками. `POST /api/teachers/{id}/dismiss` сохраняет преподавателя и назначения, меняя его активность. Обе операции возвращают `204`. После этого новые оценки отчисленному студенту и назначения уволенному преподавателю отклоняются с `400`.

При `DELETE` у студента предварительно удаляются оценки, у преподавателя — связанные занятия и назначения. Эти действия выполняются в транзакции. Остальные связанные справочники защищены внешними ключами: например, нельзя удалить группу, пока на неё ссылаются студенты или назначения.

Например, удаление студента удалило и его оценку:

```text
DELETE /api/students/1
HTTP 204 No Content
```

```text
GET /api/grades/1
HTTP 404 Not Found
{
  "error": "запись не найдена"
}
```

Ведомость строится по текущему `group_id` студента и включает сохранённых неактивных студентов этой группы. Отдельная история переводов между группами в модели не хранится.
