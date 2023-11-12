package main

import (
	"database/sql"
	"htmx-reddit/internal/db"
	"htmx-reddit/internal/service"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"

	_ "github.com/mattn/go-sqlite3"
)

const (
	dbFile = "todo.db"
)

func main() {
	// dbConn connection
	dbConn, err := sql.Open("sqlite3", dbFile)
	if err != nil {
		log.Fatal("opening db", "error", err)
	}
	db.Setup(dbConn)

	sess := scs.New()
	sess.Lifetime = 24 * time.Hour
	sess.Cookie.Secure = true

	router := httprouter.New()

	// setup routes
	routes(
		router,
		sess,
		service.NewPost(db.NewPostStore(dbConn), sess),
		service.NewComment(db.NewCommentStore(dbConn), sess),
		service.NewUser(db.NewUserStore(dbConn), sess),
	)

	log.Fatal(
		"ended execution",
		"error", http.ListenAndServe("localhost:4001", router),
	)
}
