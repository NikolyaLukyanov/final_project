package api

import (
	"net/http"
	"time"

	"go1f/pkg/db"
)

func NewApp(storage *db.Storage) *App {
	return &App{Storage: storage}
}

func (a *App) Init() {
	http.HandleFunc("/api/signin", signinHandler)
	http.HandleFunc("/api/nextdate", nextDateHandler)

	http.HandleFunc("/api/task", auth(a.taskHandler))
	http.HandleFunc("/api/tasks", auth(a.tasksHandler))
	http.HandleFunc("/api/task/done", auth(a.doneHandler))

	http.Handle("/", http.FileServer(http.Dir("web")))
}

func nextDateHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	nowStr := r.FormValue("now")
	dateStr := r.FormValue("date")
	repeat := r.FormValue("repeat")

	var now time.Time
	if nowStr == "" {
		now = time.Now()
	} else {
		var err error
		now, err = time.Parse(dateFormat, nowStr)
		if err != nil {
			http.Error(w, "Invalid 'now' date format", http.StatusBadRequest)
			return
		}
	}

	nextDate, err := NextDate(now, dateStr, repeat)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(nextDate))
}
