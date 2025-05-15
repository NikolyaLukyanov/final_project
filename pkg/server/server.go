package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"

	"go1f/pkg/api"
)

func Run(app *api.App) error {
	port := 7540
	if envPort := os.Getenv("TODO_PORT"); envPort != "" {
		if p, err := strconv.Atoi(envPort); err == nil {
			port = p
		}
	}

	app.Init() // Инициализация маршрутов

	fmt.Printf("Сервер запущен на порту %d\n", port)
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
