package main

import (
	mw "htmx-reddit/internal/middleware"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/web/components"
	"htmx-reddit/internal/web/pages"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/julienschmidt/httprouter"
)

// TODO route consts that are passed to the template?
func routes(
	r *httprouter.Router,
	sess *scs.SessionManager,
	postService service.Post,
	commentService service.Comment,
	userService service.User,
) {
	// CSS
	r.ServeFiles("/static/css/*filepath", http.Dir("./css/"))

	mws := mw.Join(sess.LoadAndSave, mw.AuthenticateMiddleware(sess, userService))

	// pages
	pages := pages.NewHandler(commentService, postService, userService)
	r.Handler("GET", "/", mws(pages.Home))
	r.Handler("GET", "/posts", mws(pages.AllPosts))
	r.Handler("GET", "/posts/new", mws(pages.NewPost))
	r.Handler("GET", "/post/:id", mws(pages.Post))
	r.Handler("GET", "/users/new", mws(pages.NewUser))

	// partial htmx data
	components := components.NewHandler(commentService, postService, userService)

	// Comment Partials
	r.Handler("POST", "/comment/add", mws(components.Comment.Add))
	r.Handler("DELETE", "/comment/delete/:id", mws(components.Comment.Delete))
	r.Handler("POST", "/comment/reply/:id", mws(components.Comment.Reply))
	r.Handler("POST", "/reply/show", mws(components.Comment.ShowReply))
	r.Handler("POST", "/reply/hide", mws(components.Comment.HideReply))

	// Post Partials
	r.Handler("POST", "/post/add", mws(components.Post.Add))
	r.Handler("DELETE", "/post/delete/:id", mws(components.Post.Delete))

	// User Partials
	r.Handler("POST", "/user/add", mws(components.User.Add))
}
