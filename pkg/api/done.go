package api

import (
	"net/http"
	"time"
)

func (a *App) doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		http.Error(w, `{"error":"Не указан ID задачи"}`, http.StatusBadRequest)
		return
	}

	task, err := a.Storage.GetTask(id)
	if err != nil {
		http.Error(w, `{"error":"Задача не найдена"}`, http.StatusNotFound)
		return
	}

	if task.Repeat == "" {
		if err := a.Storage.DeleteTask(id); err != nil {
			http.Error(w, `{"error":"Ошибка удаления задачи"}`, http.StatusInternalServerError)
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		http.Error(w, `{"error":"Ошибка расчёта следующей даты"}`, http.StatusBadRequest)
		return
	}

	if err := a.Storage.UpdateDate(next, id); err != nil {
		http.Error(w, `{"error":"Ошибка обновления даты"}`, http.StatusInternalServerError)
		return
	}

	writeJSON(w, map[string]string{})
}
