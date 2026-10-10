package analytics

type Lesson struct {
	ScheduleID        int64  `json:"schedule_id"`
	LessonNumber      int    `json:"lesson_number"`
	SubjectID         int64  `json:"subject_id"`
	SubjectName       string `json:"subject_name"`
	TeacherID         int64  `json:"teacher_id"`
	TeacherFirstName  string `json:"teacher_first_name"`
	TeacherLastName   string `json:"teacher_last_name"`
	ClassroomID       int64  `json:"classroom_id"`
	ClassroomNumber   string `json:"classroom_number"`
	ClassroomBuilding string `json:"classroom_building"`
}

type TeacherResult struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type GroupResult struct {
	ID     int64  `json:"id"`
	Name   string `json:"name"`
	Course int    `json:"course"`
}

type CourseCount struct {
	Course       int   `json:"course"`
	StudentCount int64 `json:"student_count"`
}
