package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go1f/pkg/db"
)

func taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		addTaskHandler(w, r)
	case http.MethodPut:
		editTaskHandler(w, r)
	case http.MethodGet:
		getTaskHandler(w, r)
	case http.MethodDelete:
		deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	// Чтение JSON
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка чтения JSON: " + err.Error()})
		return
	}

	// Проверка Title
	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	// Проверка и обработка даты
	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Добавление задачи в базу
	id, err := db.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка добавления в БД: " + err.Error()})
		return
	}

	// Успешный ответ
	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJSON(w, task)
}

func editTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	if err := json.NewDecoder(r.Body).Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка чтения JSON"})
		return
	}

	if task.ID == "" {
		writeJSON(w, map[string]string{"error": "Не указан ID задачи"})
		return
	}

	if task.Title == "" {
		writeJSON(w, map[string]string{"error": "Не указан заголовок задачи"})
		return
	}

	if err := checkDate(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	if err := db.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{}) // пустой JSON — успешное обновление
}

func deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан ID задачи"})
		return
	}

	if err := db.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка удаления задачи"})
		return
	}

	writeJSON(w, map[string]string{})
}

func checkDate(task *db.Task) error {
	now := time.Now()

	// Если дата пустая — подставляем сегодня
	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	parsedDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("Дата указана в неверном формате")
	}

	if task.Repeat != "" {
		// Проверка правила повторения
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("Некорректное правило повторения")
		}
		// Если дата в прошлом, подставляем следующую
		if now.After(parsedDate) {
			task.Date = next
		}
	} else {
		// Если правило пустое и дата в прошлом
		if now.After(parsedDate) {
			task.Date = now.Format(dateFormat)
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(data)
}
