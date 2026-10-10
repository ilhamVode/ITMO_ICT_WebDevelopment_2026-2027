-- 1. Учебные группы
CREATE TABLE IF NOT EXISTS study_groups (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(50) NOT NULL UNIQUE,
    course INT NOT NULL CHECK (course > 0)
);

-- 2. Кабинеты
CREATE TABLE IF NOT EXISTS classrooms (
    id BIGSERIAL PRIMARY KEY,
    number VARCHAR(20) NOT NULL,
    building VARCHAR(100) NOT NULL,

    UNIQUE (number, building)
);

-- 3. Дисциплины
CREATE TABLE IF NOT EXISTS subjects (
    id BIGSERIAL PRIMARY KEY,
    name VARCHAR(150) NOT NULL UNIQUE
);

-- 4. Семестры
CREATE TABLE IF NOT EXISTS semesters (
    id BIGSERIAL PRIMARY KEY,
    academic_year VARCHAR(9) NOT NULL,
    number INT NOT NULL CHECK (number IN (1, 2)),

    UNIQUE (academic_year, number)
);

-- 5. Студенты
CREATE TABLE IF NOT EXISTS students (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    birth_date DATE,
    group_id BIGINT NOT NULL,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    FOREIGN KEY (group_id)
        REFERENCES study_groups(id)
        ON DELETE RESTRICT
);

-- 6. Преподаватели
CREATE TABLE IF NOT EXISTS teachers (
    id BIGSERIAL PRIMARY KEY,
    first_name VARCHAR(100) NOT NULL,
    last_name VARCHAR(100) NOT NULL,
    middle_name VARCHAR(100),
    fixed_classroom_id BIGINT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,

    FOREIGN KEY (fixed_classroom_id)
        REFERENCES classrooms(id)
        ON DELETE SET NULL
);

-- 7. Назначения преподавателей
CREATE TABLE IF NOT EXISTS teaching_assignments (
    id BIGSERIAL PRIMARY KEY,
    teacher_id BIGINT NOT NULL,
    subject_id BIGINT NOT NULL,
    group_id BIGINT NOT NULL,
    semester_id BIGINT NOT NULL,

    FOREIGN KEY (teacher_id)
        REFERENCES teachers(id),

    FOREIGN KEY (subject_id)
        REFERENCES subjects(id),

    FOREIGN KEY (group_id)
        REFERENCES study_groups(id),

    FOREIGN KEY (semester_id)
        REFERENCES semesters(id),

    UNIQUE (teacher_id, subject_id, group_id, semester_id)
);

-- 8. Расписание
CREATE TABLE IF NOT EXISTS schedules (
    id BIGSERIAL PRIMARY KEY,
    assignment_id BIGINT NOT NULL,
    classroom_id BIGINT NOT NULL,
    weekday INT NOT NULL CHECK (weekday BETWEEN 1 AND 7),
    lesson_number INT NOT NULL CHECK (lesson_number > 0),

    FOREIGN KEY (assignment_id)
        REFERENCES teaching_assignments(id),

    FOREIGN KEY (classroom_id)
        REFERENCES classrooms(id),

    UNIQUE (assignment_id, weekday, lesson_number)
);

-- 9. Семестровые оценки
CREATE TABLE IF NOT EXISTS grades (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    subject_id BIGINT NOT NULL,
    semester_id BIGINT NOT NULL,
    value INT NOT NULL CHECK (value BETWEEN 2 AND 5),

    FOREIGN KEY (student_id)
        REFERENCES students(id),

    FOREIGN KEY (subject_id)
        REFERENCES subjects(id),

    FOREIGN KEY (semester_id)
        REFERENCES semesters(id),

    UNIQUE (student_id, subject_id, semester_id)
);

-- 10. Пользователи API. Пароли хранятся в виде bcrypt-хеша.
CREATE TABLE IF NOT EXISTS api_users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(50) NOT NULL UNIQUE,
    password_hash VARCHAR(100) NOT NULL
);

-- 11. Токены авторизации. В key хранится SHA-256, а не исходный токен.
CREATE TABLE IF NOT EXISTS api_tokens (
    key VARCHAR(64) PRIMARY KEY,
    user_id BIGINT NOT NULL REFERENCES api_users(id) ON DELETE CASCADE,
    expires_at TIMESTAMPTZ NOT NULL
);

CREATE INDEX IF NOT EXISTS students_group_idx ON students(group_id);
CREATE INDEX IF NOT EXISTS teaching_group_semester_idx ON teaching_assignments(group_id, semester_id);
CREATE INDEX IF NOT EXISTS schedules_slot_idx ON schedules(weekday, lesson_number);
CREATE INDEX IF NOT EXISTS grades_semester_idx ON grades(semester_id);
