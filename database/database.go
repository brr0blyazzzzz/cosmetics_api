package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

// инициализация и проверка подключения к базе данных
func InitDB() error {
	var err error
	DB, err = sql.Open("sqlite3", "./cosmetics.db")
	if err != nil {
		return err
	}

	err = DB.Ping()
	if err != nil {
		return err
	}

	log.Println("Подключение к базе косметических продуктов успешно :)")
	return nil
}
