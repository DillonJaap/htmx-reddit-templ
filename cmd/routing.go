package main

import (
	commentDB "htmx-reddit/internal/models/comment"
	postDB "htmx-reddit/internal/models/post"
	userDB "htmx-reddit/internal/models/user"
	handlers "htmx-reddit/internal/web"
	"htmx-reddit/internal/web/pages"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO route consts that are passed to the template?
func routes(
	r *httprouter.Router,
	postModel postDB.Model,
	cmntModel commentDB.Model,
	userModel userDB.Model,
) {
	// CSS
	r.ServeFiles("/css/*filepath", http.Dir("./ui/css/"))

	// pages
	pages := pages.New(cmntModel, postModel, userModel)
	r.GET("/", pages.Home)
	r.GET("/posts", pages.AllPosts)
	r.GET("/posts/new", pages.NewPost)
	r.GET("/post/:id", pages.Post)
	r.GET("/users/new", pages.NewUser)

	// partial htmx data
	handler := handlers.New(cmntModel, postModel, userModel)
	r.POST("/comment/add", handler.Comment.Add)
	r.DELETE("/comment/delete/:id", handler.Comment.Delete)
	r.POST("/comment/reply/:id", handler.Comment.Reply)
	r.POST("/reply/show", handler.Comment.ShowReply)
	r.POST("/reply/hide", handler.Comment.HideReply)

	r.POST("/post/add", handler.Post.Add)
	r.DELETE("/post/delete/:id", handler.Post.Delete)

	r.POST("/user/add", handler.User.Add)
}
