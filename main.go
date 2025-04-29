package main

import (
	"log"
	"os"

	"go1f/pkg/db"
	"go1f/pkg/server"

	"github.com/joho/godotenv"
)

func main() {
	dbFile := "scheduler.db"
	godotenv.Load()
	if envPath := os.Getenv("TODO_DBFILE"); envPath != "" {
		dbFile = envPath
	}

	if err := db.Init(dbFile); err != nil {
		log.Fatalf("Ошибка инициализации базы данных: %v", err)
	}

	if err := server.Run(); err != nil {
		log.Fatalf("Ошибка запуска сервера: %v", err)
	}
}
