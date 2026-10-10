package httpapi

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5/pgconn"
	"gorm.io/gorm"
)

// Отправляем JSON клиенту
func WriteJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(data); err != nil {
		log.Printf("Ошибка отправки JSON: %v", err)
	}
}

// Отправляем ошибку клиенту
func WriteError(w http.ResponseWriter, status int, message string) {
	WriteJSON(w, status, map[string]string{
		"error": message,
	})
}

// Читаем JSON из тела запроса
func ReadJSON(w http.ResponseWriter, r *http.Request, dst any) error {
	decoder := json.NewDecoder(http.MaxBytesReader(w, r.Body, 1<<20))
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(dst); err != nil {
		return err
	}

	var extra any
	if err := decoder.Decode(&extra); err != io.EOF {
		return errors.New("ожидался только один JSON-объект")
	}

	return nil
}

// Получаем ID из адреса /api/groups/{id}
func PathID(r *http.Request) (int64, error) {
	id, err := strconv.ParseInt(r.PathValue("id"), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.New("некорректный ID")
	}

	return id, nil
}

// Преобразуем ошибки сервиса и БД в HTTP-ответы
func HandleError(w http.ResponseWriter, err error, invalid error) {
	switch {
	case errors.Is(err, invalid):
		WriteError(w, http.StatusBadRequest, err.Error())
		return

	case errors.Is(err, gorm.ErrRecordNotFound):
		WriteError(w, http.StatusNotFound, "запись не найдена")
		return
	}

	var pgErr *pgconn.PgError

	if errors.As(err, &pgErr) {
		switch pgErr.Code {
		case "23P01":
			WriteError(w, http.StatusConflict, pgErr.Message)
			return
		case "23505":
			WriteError(w, http.StatusConflict, "запись уже существует")
			return

		case "23503":
			WriteError(w, http.StatusConflict, "нарушена связь с другой таблицей")
			return

		case "23514":
			WriteError(w, http.StatusBadRequest, pgErr.Message)
			return
		case "23502", "22001", "22003":
			WriteError(w, http.StatusBadRequest, "нарушены ограничения базы данных")
			return
		}
	}

	log.Printf("Ошибка API: %v", err)
	WriteError(w, http.StatusInternalServerError, "внутренняя ошибка сервера")
}
