package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

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

type Storage struct {
	DB *sql.DB
}

func (s *Storage) Init(dbFile string) error {
	_, err := os.Stat(dbFile)
	install := os.IsNotExist(err)

	db, err := sql.Open("sqlite", dbFile)
	if err != nil {
		return fmt.Errorf("ошибка открытия базы данных: %w", err)
	}

	if err = db.Ping(); err != nil {
		return fmt.Errorf("ошибка соединения с базой: %w", err)
	}

	s.DB = db

	if install {
		if _, err = db.Exec(schema); err != nil {
			return fmt.Errorf("ошибка создания схемы: %w", err)
		}
	}

	return nil
}

func (s *Storage) Close() error {
	if s.DB != nil {
		return s.DB.Close()
	}
	return nil
}
