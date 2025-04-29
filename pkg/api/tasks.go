package api

import (
	"net/http"
	"strings"

	"go1f/pkg/db"
)

type TasksResp struct {
	Tasks []*db.Task `json:"tasks"`
}

func tasksHandler(w http.ResponseWriter, r *http.Request) {
	search := strings.TrimSpace(r.URL.Query().Get("search"))

	tasks, err := db.Tasks(50, search)
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
