package pages

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/helpers"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/templ"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Home     http.HandlerFunc
	Post     http.HandlerFunc
	AllPosts http.HandlerFunc
	NewPost  http.HandlerFunc
	SignUp   http.HandlerFunc
	Login    http.HandlerFunc
}

func NewHandler(
	c service.Comment,
	p service.Post,
	u service.User,
	s *scs.SessionManager,
) *Handler {
	return &Handler{
		Home:     allPosts(p, s),
		AllPosts: allPosts(p, s),
		Post:     postPage(p, c, s),
		NewPost:  newPost(s),
		SignUp:   signup(s),
		Login:    login(),
	}
}

func allPosts(ps service.Post, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		postList, err := ps.GetAll()
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}

		templ.AllPosts(
			postList,
			helpers.IsAuthenticated(req),
			sess.GetString(req.Context(), "username"),
		).Render(req.Context(), w)
	}
}

func newPost(sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		templ.NewPost(
			helpers.IsAuthenticated(req),
			sess.GetString(req.Context(), "username"),
		).Render(req.Context(), w)
	}
}

func signup(sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var pageData struct {
			ShowPassErr bool
		}
		pageData.ShowPassErr = false
		templ.SignUp(
			helpers.IsAuthenticated(req),
			sess.GetString(req.Context(), "username"),
		).Render(req.Context(), w)
	}
}

func login() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var pageData struct {
			ShowPassErr bool
		}
		pageData.ShowPassErr = false
		templ.Login().Render(req.Context(), w)
	}
}

func postPage(postService service.Post, commentService service.Comment, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p := httprouter.ParamsFromContext(req.Context())
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("failed to convert to int", "error", err)
		}

		post, err := postService.Get(id)
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}
		log.Printf("%+v", post)

		comments := commentService.GetByPostID(post.ID)

		templ.Post(
			post,
			comments,
			helpers.IsAuthenticated(req),
			sess.GetString(req.Context(), "username"),
		).Render(req.Context(), w)
	}
}
