package main

import (
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/web/components"
	"htmx-reddit/internal/web/pages"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

// TODO route consts that are passed to the template?
func routes(
	r *httprouter.Router,
	postService service.Post,
	commentService service.Comment,
	userService service.User,
) {
	// CSS
	r.ServeFiles("/css/*filepath", http.Dir("./ui/css/"))

	// pages
	pages := pages.NewHandler(commentService, postService, userService)
	r.GET("/", pages.Home)
	r.GET("/posts", pages.AllPosts)
	r.GET("/posts/new", pages.NewPost)
	r.GET("/post/:id", pages.Post)
	r.GET("/users/new", pages.NewUser)

	// partial htmx data
	components := components.NewHandler(commentService, postService, userService)
	r.POST("/comment/add", components.Comment.Add)
	r.DELETE("/comment/delete/:id", components.Comment.Delete)
	r.POST("/comment/reply/:id", components.Comment.Reply)
	r.POST("/reply/show", components.Comment.ShowReply)
	r.POST("/reply/hide", components.Comment.HideReply)

	r.POST("/post/add", components.Post.Add)
	r.DELETE("/post/delete/:id", components.Post.Delete)

	r.POST("/user/add", components.User.Add)
}
