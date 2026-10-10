-- Проверяем пересечения расписания на уровне БД, включая одновременные запросы.
CREATE OR REPLACE FUNCTION check_schedule() RETURNS TRIGGER AS $$
DECLARE
    assignment teaching_assignments%ROWTYPE;
BEGIN
    -- Общий замок для изменений расписания и назначений в текущей транзакции.
    PERFORM pg_advisory_xact_lock(312012);
    SELECT * INTO assignment FROM teaching_assignments WHERE id = NEW.assignment_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'назначение не найдено' USING ERRCODE = '23503';
    END IF;
    IF NOT EXISTS (SELECT 1 FROM teachers WHERE id = assignment.teacher_id AND is_active) THEN
        RAISE EXCEPTION 'уволенному преподавателю нельзя назначать занятия' USING ERRCODE = '23514';
    END IF;
    IF EXISTS (
        SELECT 1 FROM schedules sc
        JOIN teaching_assignments ta ON ta.id = sc.assignment_id
        WHERE sc.id <> COALESCE(NEW.id, 0)
          AND sc.weekday = NEW.weekday AND sc.lesson_number = NEW.lesson_number
          AND ta.semester_id = assignment.semester_id
          AND (ta.group_id = assignment.group_id OR ta.teacher_id = assignment.teacher_id
               OR sc.classroom_id = NEW.classroom_id)
    ) THEN
        RAISE EXCEPTION 'в эту пару заняты группа, преподаватель или кабинет' USING ERRCODE = '23P01';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER schedule_check
BEFORE INSERT OR UPDATE ON schedules FOR EACH ROW EXECUTE FUNCTION check_schedule();

-- Изменение назначения тоже не должно создавать пересечения расписания.
CREATE OR REPLACE FUNCTION check_teaching_assignment() RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_advisory_xact_lock(312012);
    IF NOT EXISTS (SELECT 1 FROM teachers WHERE id = NEW.teacher_id AND is_active) THEN
        RAISE EXCEPTION 'назначить можно только работающего преподавателя' USING ERRCODE = '23514';
    END IF;
    IF EXISTS (
        SELECT 1 FROM schedules own_sc
        JOIN schedules other_sc ON other_sc.weekday = own_sc.weekday
            AND other_sc.lesson_number = own_sc.lesson_number AND other_sc.id <> own_sc.id
        JOIN teaching_assignments other_ta ON other_ta.id = other_sc.assignment_id
        WHERE own_sc.assignment_id = NEW.id AND other_ta.id <> NEW.id
          AND other_ta.semester_id = NEW.semester_id
          AND (other_ta.group_id = NEW.group_id OR other_ta.teacher_id = NEW.teacher_id
               OR other_sc.classroom_id = own_sc.classroom_id)
    ) THEN
        RAISE EXCEPTION 'изменение назначения создаёт конфликт расписания' USING ERRCODE = '23P01';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER teaching_assignment_check
BEFORE INSERT OR UPDATE ON teaching_assignments
FOR EACH ROW EXECUTE FUNCTION check_teaching_assignment();

-- Оценку выставляем обучающемуся студенту по дисциплине его группы за семестр.
CREATE OR REPLACE FUNCTION check_grade() RETURNS TRIGGER AS $$
DECLARE
    student_group BIGINT;
    student_active BOOLEAN;
BEGIN
    SELECT group_id, is_active INTO student_group, student_active FROM students WHERE id = NEW.student_id;
    IF NOT FOUND THEN
        RAISE EXCEPTION 'студент не найден' USING ERRCODE = '23503';
    END IF;
    IF NOT student_active THEN
        RAISE EXCEPTION 'отчисленному студенту нельзя выставлять оценки' USING ERRCODE = '23514';
    END IF;
    IF NOT EXISTS (
        SELECT 1 FROM teaching_assignments WHERE group_id = student_group
          AND subject_id = NEW.subject_id AND semester_id = NEW.semester_id
    ) THEN
        RAISE EXCEPTION 'дисциплина не назначена группе студента в этом семестре' USING ERRCODE = '23514';
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER grade_check
BEFORE INSERT OR UPDATE ON grades FOR EACH ROW EXECUTE FUNCTION check_grade();
