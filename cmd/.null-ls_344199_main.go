package main

import (
	"database/sql"
	"htmx-reddit/internal/db"
	"htmx-reddit/internal/service"
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

	// sqlDB connection
	sqlDB, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("opening db", "error", err)
	}

	// setup routes
	routes(
		router,
		service.NewPost(db.NewPostStore(sqlDB)),
		service.NewComment(db.NewCommentStore(sqlDB)),
		service.NewUser(db.NewUserStore(sqlDB)),
	)

	log.Fatal(
		"ended execution",
		"error", http.ListenAndServe("localhost:4001", router),
	)
}
