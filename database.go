package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func createDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "db/data.db?_foreign_keys=ON&_journal_mode=WAL")

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

		create table if not exists rooms (
			name text,
			body text,
			links text,
			image blob,
			style text,
			user_id integer not null,
			foreign key (user_id) references users(id) on delete cascade
		);

		create table if not exists logos (
			name text,
			image blob
		);
	`

	_, err = db.Exec(sqlStatement)

	if err != nil {
		return db, err
	}

	return db, nil
}
