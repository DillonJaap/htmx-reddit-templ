package main

import (
	"database/sql"
	"htmx-reddit/internal/models/comment"
	"htmx-reddit/internal/models/post"
	"htmx-reddit/internal/models/user"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "todo.db"
)

func main() {
	router := httprouter.New()

	// db connection
	db, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("opening db", "error", err)
	}

	// setup routes
	routes(
		router,
		post.New(db),
		comment.New(db),
		user.New(db),
	)

	log.Fatal(
		"ended execution",
		"error", http.ListenAndServe("localhost:4001", router),
	)
}
