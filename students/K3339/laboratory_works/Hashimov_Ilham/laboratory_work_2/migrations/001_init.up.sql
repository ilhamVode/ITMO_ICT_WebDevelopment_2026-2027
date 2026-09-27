CREATE TABLE users (
    id BIGSERIAL PRIMARY KEY,
    username VARCHAR(100) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    is_admin BOOLEAN NOT NULL DEFAULT FALSE,
    class_name VARCHAR(50)
);

CREATE TABLE homeworks (
    id BIGSERIAL PRIMARY KEY,
    subject VARCHAR(100) NOT NULL,
    teacher_id BIGINT NOT NULL,
    issue_date TIMESTAMP NOT NULL,
    start_date TIMESTAMP NOT NULL,
    end_date TIMESTAMP NOT NULL,
    text TEXT NOT NULL,
    penalty TEXT,

    FOREIGN KEY (teacher_id) REFERENCES users(id)
);

CREATE TABLE submissions (
    id BIGSERIAL PRIMARY KEY,
    student_id BIGINT NOT NULL,
    homework_id BIGINT NOT NULL,
    text TEXT NOT NULL,
    submitted_at TIMESTAMP NOT NULL,
    grade DOUBLE PRECISION,

    FOREIGN KEY (student_id) REFERENCES users(id),
    FOREIGN KEY (homework_id) REFERENCES homeworks(id)
);