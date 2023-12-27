package main

import (
	mw "htmx-reddit/internal/middleware"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/web/pages"
	partials "htmx-reddit/internal/web/partials"
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

	// setup middlewares
	defaultMW := mw.Join(
		sess.LoadAndSave,
		mw.Authenticate(sess, userService),
	)

	authMW := mw.Join(
		defaultMW,
		mw.RequireAuthentication(),
	)

	// pages
	pages := pages.NewHandler(commentService, postService, userService, sess)
	r.Handler("GET", "/", defaultMW(pages.Home))
	r.Handler("GET", "/posts", defaultMW(pages.AllPosts))
	r.Handler("GET", "/posts/new", defaultMW(pages.NewPost))
	r.Handler("GET", "/post/:id", defaultMW(pages.Post))
	r.Handler("GET", "/users/new", defaultMW(pages.SignUp))
	r.Handler("GET", "/users/login", defaultMW(pages.Login))

	// partial htmx data
	partials := partials.NewHandler(sess, commentService, postService, userService)

	// Comment Partials
	r.Handler("POST", "/comment/add", authMW(partials.Comment.Add))
	r.Handler("DELETE", "/comment/delete/:id", authMW(partials.Comment.Delete))
	r.Handler("POST", "/comment/reply/:id", authMW(partials.Comment.Reply))

	r.Handler("POST", "/reply/show", defaultMW(partials.Comment.ShowReply))
	r.Handler("POST", "/reply/hide", defaultMW(partials.Comment.HideReply))

	// Post Partials
	r.Handler("POST", "/post/add", authMW(partials.Post.Add))
	r.Handler("DELETE", "/post/delete/:id", authMW(partials.Post.Delete))

	// User Partials
	r.Handler("POST", "/user/add", defaultMW(partials.User.Add))
	r.Handler("POST", "/user/login", defaultMW(partials.User.Login))
	r.Handler("POST", "/user/logout", authMW(partials.User.Logout))
}
