package db

import (
	"database/sql"

	"github.com/charmbracelet/log"
)

func Setup(conn *sql.DB) {
	_, err := conn.Exec(`
	CREATE TABLE IF NOT EXISTS user (
		id              INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		name            TEXT NOT NULL,
		hashed_password TEXT NOT NULL,
		created_time    DATETIME NOT NULL
	);
	CREATE TABLE IF NOT EXISTS post (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		title TEXT NOT NULL,
		body TEXT NOT NULL,
		created_time DATETIME NOT NULL,
		owner TEXT NOT NULL,
		owner_id  INTEGER,
		FOREIGN KEY (owner_id) REFERENCES user (id)
	);
	CREATE TABLE IF NOT EXISTS comment (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		description TEXT NOT NULL,
		created_time DATETIME NOT NULL,
		owner TEXT NOT NULL,
		parent_id INTEGER,
		post_id INTEGER,
		owner_id  INTEGER,
		FOREIGN KEY (parent_id) REFERENCES comment (id),
		FOREIGN KEY (post_id) REFERENCES post (id),
		FOREIGN KEY (owner_id) REFERENCES user (id)
	);
	`)
	if err != nil {
		log.Fatalf("error creating table: %s", err)
	}

}
