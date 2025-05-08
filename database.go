package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "users.db")

	if err != nil {
		return db, err
	}

	sqlStatement := `
		create table if not exists users (
			id integer not null primary key autoincrement,
			name text,
			email text,
			hash text,
			session text,
			csrf text
		);
	`

	_, err = db.Exec(sqlStatement)

	if err != nil {
		return db, err
	}

	return db, nil
}
