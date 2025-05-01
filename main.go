package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	"go1f/pkg/api"
	"go1f/pkg/db"
)

const defaultDBFile = "scheduler.db"
const defaultPort = "7540"

func main() {
	_ = godotenv.Load()
	log.Println("TODO_PASSWORD:", os.Getenv("TODO_PASSWORD"))

	dbFile := os.Getenv("TODO_DBFILE")
	if dbFile == "" {
		dbFile = defaultDBFile
	}

	storage := &db.Storage{}
	if err := storage.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}
	defer storage.Close()

	app := api.NewApp(storage)
	app.Init()

	port := os.Getenv("TODO_PORT")
	if port == "" {
		port = defaultPort
	}

	log.Printf("Сервер запущен на порту %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
