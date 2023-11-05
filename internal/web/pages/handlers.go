package pages

import (
	"context"
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/templ"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Home     http.HandlerFunc
	Post     http.HandlerFunc
	AllPosts http.HandlerFunc
	NewPost  http.HandlerFunc
	NewUser  http.HandlerFunc
}

func NewHandler(
	c service.Comment,
	p service.Post,
	u service.User,
) *Handler {
	return &Handler{
		Home:     allPosts(p),
		AllPosts: allPosts(p),
		Post:     postPage(p, c),
		NewPost:  newPost(),
		NewUser:  newUser(),
	}
}

func allPosts(ps service.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		postList, err := ps.GetAll()
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}

		templ.AllPosts(postList).Render(context.TODO(), w)
	}
}

func newPost() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		templ.NewPost().Render(context.TODO(), w)
	}
}

func newUser() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		var pageData struct {
			ShowPassErr bool
		}
		pageData.ShowPassErr = false
		templ.NewUser().Render(context.TODO(), w)
	}
}

func postPage(postService service.Post, commentService service.Comment) http.HandlerFunc {
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

		comments := commentService.GetByPostID(post.ID)
		templ.Post(post, comments).Render(context.TODO(), w)
	}
}
