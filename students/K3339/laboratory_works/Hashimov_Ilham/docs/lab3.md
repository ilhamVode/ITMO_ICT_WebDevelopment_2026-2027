# Лабораторная работа №3

## Вариант 12. Учебная часть колледжа

В третьей работе реализован REST API для учебной части колледжа. В нём хранятся группы, студенты, преподаватели, дисциплины, кабинеты, назначения, расписание и семестровые оценки. Работа на Go, с GORM и PostgreSQL.

В отличие от второй лабы, обработчики возвращают JSON. Ответ получает клиент: `curl`, Postman или программа, которая вызывает API. В консоли сервера остаются сообщения о запуске и ошибках.

Ниже приведены ответы, полученные при проверке обработчиков с настоящей PostgreSQL на тестовых данных. Они оформлены как запрос, HTTP-статус и тело ответа, без служебных заголовков. Значение токена заменено на `<token>`. ID относятся к тестовой базе, в другой базе они могут отличаться.

---

## Task 1. Создание и запуск проекта

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

Перед запуском нужна созданная база PostgreSQL. Локальный `.env` заполняется по `.env.example`:

```dotenv
DATABASE_DSN=host=localhost port=5432 user=postgres password=change_me dbname=college_lab3 sslmode=disable
SERVER_ADDR=:8080
```

Здесь приведён пример настройки, пароль нужно заменить. `.env` и результаты сборки исключены из Git через `.gitignore`.

Запуск из папки третьей лабораторной:

```powershell
go mod download
go run ./cmd/server
```

В `main.go` подключение создаётся через `database.NewPostgres(dsn)`. Затем `database.InitSchema(db)` применяет оба SQL-файла в одной транзакции. Таблицы создаются через `CREATE TABLE IF NOT EXISTS`, поэтому повторный запуск не удаляет записи. Функции и триггеры обновляются через `CREATE OR REPLACE`.

В отличие от Django, здесь нет команд `makemigrations` и `migrate`: схема написана в SQL и применяется кодом приложения. GORM используется для работы с моделями и запросами.

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

Связь преподавателей и дисциплин является m2m: преподаватель ведёт несколько дисциплин, одну дисциплину ведут разные преподаватели. Промежуточная модель `TeachingAssignment` дополнительно хранит группу и семестр. По назначению это соответствует `through` в Django.

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

Запросы разобраны на моделях выбранного варианта. Создание, получение и фильтрация через GORM имеют вид:

```go
err := db.WithContext(ctx).Create(&student).Error
err := db.WithContext(ctx).First(&student, id).Error
err := db.WithContext(ctx).Order("id ASC").Find(&students).Error
err := db.WithContext(ctx).Where("group_id = ?", groupID).Find(&students).Error
```

Это соответствует операциям `objects.create`, `get`, `all` и `filter` в Django. `&student` передаёт адрес структуры, чтобы ORM мог заполнить её, включая ID созданной записи. Ошибка возвращается значением `error`, а не исключением.

Вместо `related_name` обратная связь описана полем структуры и тегом внешнего ключа:

```go
Students []student.Student `gorm:"foreignKey:GroupID;references:ID" json:"students,omitempty"`
```

`Students` не является отдельной колонкой. Для загрузки связанных студентов используется `Preload("Students")`, а для дисциплин преподавателя — `Preload("Assignments.Subject")`.

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

Вместо `APIView` и generic-классов DRF используются методы структур. Вместо отдельного класса сериализатора JSON читается в структуру запроса, проверяется сервисом и преобразуется в модель. Ответ сериализуется через стандартный `encoding/json`:

```go
var req gradeRequest
err := httpapi.ReadJSON(w, r, &req)
grade := req.ToModel()
err = h.service.Create(r.Context(), &grade)
httpapi.WriteJSON(w, http.StatusCreated, grade)
```

Это отдельные строки основного пути обработки; в полном обработчике после каждого действия проверяется `err`. Неизвестные поля и лишний JSON-объект в теле запроса отклоняются. Даты рождения принимаются в формате `YYYY-MM-DD`.

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

Регистрация, вход по токену и вывод текущего пользователя соответствуют пункту 5 задания. Djoser в Go не используется: эти действия реализованы в `internal/auth`. Выход добавлен для удаления выданного токена.

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

Проверка токена подключена к учебным эндпоинтам через `RequireAuth`. Каждый запрос передаёт `Authorization: Token <token>`; также поддерживается `Bearer`. Роли администратора и обычного пользователя отдельно не вводились: зарегистрированные пользователи имеют одинаковый доступ к учебным данным. `/health`, регистрация и вход публичные.

Выход проверялся после остальных действий с учебными данными. Запрос с прежним токеном отклоняется:

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

Один HTTP-запрос возвращает группу вместе со студентами. Внутри ORM выполняет отдельные SQL-запросы для загрузки связанных записей.

### Многие ко многим: преподаватель и дисциплины

`GET /api/teachers/{id}` возвращает назначения с вложенной дисциплиной и отдельный список `subjects`. Используется `Preload("Assignments.Subject")`; список дисциплин собирается без повторов по `subject_id`.

Для компактного показа из реального ответа выбраны `id`, `first_name`, `last_name` и `subjects`:

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

Математика в назначениях этого преподавателя встречается в нескольких группах и семестрах, а в `subjects` выводится один раз. Поле `subjects` имеет тег `gorm:"-"`: это подготовленная часть ответа, а не колонка БД.

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

Ответ содержит полный список занятий с теми же полями, что и в первом запросе. Для краткости результат показан таблицей:

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

Средний по математике — 4, по физике — 4. Средний по группе считается по всем выставленным оценкам: `(3 + 5 + 4) / 3 = 4`. Это не среднее уже округлённых средних баллов студентов.

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

`DELETE` действительно удаляет запись: у студента предварительно удаляются оценки, у преподавателя — связанные занятия и назначения. Эти действия выполняются в транзакции. Остальные связанные справочники защищены внешними ключами: например, нельзя удалить группу, пока на неё ссылаются студенты или назначения.

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

---

## Проверка работы и документация

В `cmd/server/api_test.go` проверены CRUD всех девяти предметных сущностей, вложенные ответы, пять аналитических запросов, ведомость и CSV, регистрация, вход, выход, истечение токена, отчисление, увольнение и конфликты расписания. Для двух одновременных запросов на одну пару проверяется, что ответы — `201` и `409`.

`internal/report/service_test.go` отдельно проверяет средние баллы, отсутствующие оценки и пустую ведомость.

Для интеграционной проверки переменная `TEST_DATABASE_DSN` должна содержать строку подключения к доступной PostgreSQL. Например, её можно прочитать из своего `.env`:

```powershell
$settings = Get-Content .env | Where-Object { $_ -match '^DATABASE_DSN=' } | Select-Object -First 1
$env:TEST_DATABASE_DSN = ($settings -replace '^DATABASE_DSN=', '').Trim().Trim('"').Trim("'")
go test ./cmd/server ./internal/report -v -count=1
go vet ./...
go build ./...
```

Тест создаёт отдельную схему `lab3_test_...` и удаляет её после проверки. Существующие таблицы не изменяются. Без `TEST_DATABASE_DSN` интеграционные тесты пропускаются; тесты расчёта ведомости выполняются.

Вывод двух пакетов с тестами:

```text
=== RUN   TestCollegeAPI
--- PASS: TestCollegeAPI (0.65s)
=== RUN   TestCRUDResources
--- PASS: TestCRUDResources (0.35s)
PASS
ok  laboratory_work_3/cmd/server  1.032s
=== RUN   TestStatementAveragesAndMissingGrades
--- PASS: TestStatementAveragesAndMissingGrades (0.00s)
=== RUN   TestStatementWithoutGrades
--- PASS: TestStatementWithoutGrades (0.00s)
PASS
ok  laboratory_work_3/internal/report  0.024s
```

Также прошли `go vet ./...` и `go build ./...`. Они при успешном выполнении не выводят сообщений.

Отчёт находится в `docs/lab3.md` и подключён к общему `mkdocs.yml`. Для сборки документации из папки `Hashimov_Ilham`:

```powershell
python -m mkdocs build --strict
```

Серверная часть и отчёт находятся в ветке `lab3`. В этой реализации функциональность выполнена средствами Go: GORM используется вместо Django ORM, обработчики и структуры запросов — вместо DRF, `internal/auth` — вместо Djoser. Это замена стека, а не использование Python-библиотек внутри Go.
