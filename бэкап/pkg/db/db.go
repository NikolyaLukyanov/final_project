package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

var DB *sql.DB

const schema = `
CREATE TABLE scheduler (
	id INTEGER PRIMARY KEY AUTOINCREMENT,
	date CHAR(8) NOT NULL DEFAULT "",
	title VARCHAR(255) NOT NULL DEFAULT "",
	comment TEXT NOT NULL DEFAULT "",
	repeat VARCHAR(128) NOT NULL DEFAULT ""
);
CREATE INDEX date_index ON scheduler(date);
`

func Init(dbFile string) error {
	// Проверка, существует ли файл
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	// Открытие подключения к БД
	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	// Проверка соединения
	if err = db.Ping(); err != nil {
		return fmt.Errorf("ошибка соединения с базой: %w", err)
	}

	DB = db

	// Если файл создавался впервые — выполняем schema
	if install {
		if _, err = db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %w", err)
		}
	}

	return nil
}
