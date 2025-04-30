package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
)

func doneHandler(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		writeJSON(w, map[string]string{"error": "Не указан ID задачи"})
		return
	}

	task, err := db.GetTask(id)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Задача не найдена"})
		return
	}

	if task.Repeat == "" {
		// Одноразовая задача — удаляем
		if err := db.DeleteTask(id); err != nil {
			writeJSON(w, map[string]string{"error": "Ошибка удаления задачи"})
			return
		}
		writeJSON(w, map[string]string{})
		return
	}

	// Периодическая задача — обновляем дату
	now := time.Now()
	next, err := NextDate(now, task.Date, task.Repeat)
	if err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка расчёта следующей даты"})
		return
	}

	if err := db.UpdateDate(next, id); err != nil {
		writeJSON(w, map[string]string{"error": "Ошибка обновления даты"})
		return
	}

	writeJSON(w, map[string]string{})
}
