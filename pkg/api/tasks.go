package api

import (
	"net/http"
	"strings"

	"go1f/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

const maxTasksLimit = 50

func (a *App) tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	tasks, err := a.Storage.Tasks(maxTasksLimit, search)
	if err != nil {
		writeJSON(w, map[string]string{"error": err.Error()})
		return
	}

	// Гарантировать не nil, а []
	if tasks == nil {
		tasks = []*db.Task{}
	}

	writeJSON(w, TasksResp{Tasks: tasks})
}
