package report

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"laboratory_work_3/internal/httpapi"
)

type Handler struct {
	service *Service
}

func NewHandler(service *Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) RegisterRoutes(mux *http.ServeMux) {
	mux.HandleFunc("GET /api/reports/statement", h.GetStatement)
}

// GET /api/reports/statement?group_id=1&semester_id=1&format=csv
func (h *Handler) GetStatement(w http.ResponseWriter, r *http.Request) {
	groupID, err1 := strconv.ParseInt(r.URL.Query().Get("group_id"), 10, 64)
	semesterID, err2 := strconv.ParseInt(r.URL.Query().Get("semester_id"), 10, 64)
	format := r.URL.Query().Get("format")
	if err1 != nil || err2 != nil || (format != "" && format != "json" && format != "csv") {
		httpapi.WriteError(w, http.StatusBadRequest, "укажите group_id, semester_id и format=json или csv")
		return
	}
	statement, err := h.service.GetStatement(r.Context(), groupID, semesterID)
	if err != nil {
		httpapi.HandleError(w, err, ErrInvalidReport)
		return
	}
	if format == "csv" {
		writeCSV(w, statement)
		return
	}
	httpapi.WriteJSON(w, http.StatusOK, statement)
}

func averageText(value *float64) string {
	if value == nil {
		return ""
	}
	return strconv.FormatFloat(*value, 'f', 2, 64)
}

func csvText(value string) string {
	if strings.HasPrefix(value, "=") || strings.HasPrefix(value, "+") ||
		strings.HasPrefix(value, "-") || strings.HasPrefix(value, "@") ||
		strings.HasPrefix(value, "\t") || strings.HasPrefix(value, "\r") {
		return "'" + value
	}
	return value
}

func writeCSV(w http.ResponseWriter, statement Statement) {
	w.Header().Set("Content-Type", "text/csv; charset=utf-8")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=statement_%d_%d.csv", statement.Group.ID, statement.Semester.ID))
	// BOM нужен для корректного открытия русского текста в Excel.
	_, _ = w.Write([]byte{0xEF, 0xBB, 0xBF})
	writer := csv.NewWriter(w)
	writer.Comma = ';'
	rows := [][]string{
		{"Группа", csvText(statement.Group.Name)},
		{"Учебный год", statement.Semester.AcademicYear, "Семестр", strconv.Itoa(statement.Semester.Number)},
	}
	header := []string{"ID", "Фамилия", "Имя", "Обучается"}
	for _, subject := range statement.Subjects {
		header = append(header, csvText(subject.Name))
	}
	rows = append(rows, append(header, "Средний балл"))
	for _, student := range statement.Students {
		row := []string{strconv.FormatInt(student.ID, 10), csvText(student.LastName), csvText(student.FirstName), strconv.FormatBool(student.IsActive)}
		for _, grade := range student.Grades {
			value := ""
			if grade.Value != nil {
				value = strconv.Itoa(*grade.Value)
			}
			row = append(row, value)
		}
		rows = append(rows, append(row, averageText(student.Average)))
	}
	row := []string{"", "Средний балл по дисциплине", "", ""}
	for _, subject := range statement.Subjects {
		row = append(row, averageText(subject.Average))
	}
	rows = append(rows, append(row, averageText(statement.Average)))
	rows = append(rows, []string{"Количество студентов", strconv.Itoa(statement.StudentCount)},
		[]string{"Количество оценок", strconv.Itoa(statement.GradeCount)})
	if err := writer.WriteAll(rows); err != nil {
		log.Printf("Ошибка отправки CSV: %v", err)
	}
}
