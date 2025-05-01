package api

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"go1f/pkg/db"
)

type App struct {
	Storage *db.Storage
}

func (a *App) taskHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		a.addTaskHandler(w, r)
	case http.MethodPut:
		a.editTaskHandler(w, r)
	case http.MethodGet:
		a.getTaskHandler(w, r)
	case http.MethodDelete:
		a.deleteTaskHandler(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

func (a *App) addTaskHandler(w http.ResponseWriter, r *http.Request) {
	var task db.Task

	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&task); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка чтения JSON: " + err.Error()})
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

	id, err := a.Storage.AddTask(&task)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка добавления в БД: " + err.Error()})
		return
	}

	writeJSON(w, map[string]string{"id": strconv.FormatInt(id, 10)})
}

func (a *App) getTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан идентификатор"})
		return
	}

	task, err := a.Storage.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	writeJSON(w, task)
}

func (a *App) editTaskHandler(w http.ResponseWriter, r *http.Request) {
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

	if err := a.Storage.UpdateTask(&task); err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	writeJSON(w, map[string]string{})
}

func (a *App) deleteTaskHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан ID задачи"})
		return
	}

	if err := a.Storage.DeleteTask(id); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка удаления задачи"})
		return
	}

	writeJSON(w, map[string]string{})
}

func checkDate(task *db.Task) error {
	now := time.Now().Truncate(24 * time.Hour)

	if task.Date == "" {
		task.Date = now.Format(dateFormat)
	}

	parsedDate, err := time.Parse(dateFormat, task.Date)
	if err != nil {
		return errors.New("Дата указана в неверном формате")
	}

	if task.Repeat != "" {
		next, err := NextDate(now, task.Date, task.Repeat)
		if err != nil {
			return errors.New("Некорректное правило повторения")
		}

		if parsedDate.Before(now) {
			task.Date = next
		}
	} else {
		if parsedDate.Before(now) {
			task.Date = now.Format(dateFormat)
		}
	}

	return nil
}

func writeJSON(w http.ResponseWriter, data any) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if m, ok := data.(map[string]string); ok {
		if _, exists := m["error"]; exists {
			w.WriteHeader(http.StatusBadRequest)
		} else {
			w.WriteHeader(http.StatusOK)
		}
	} else {
		w.WriteHeader(http.StatusOK)
	}

	_ = json.NewEncoder(w).Encode(data)
}
