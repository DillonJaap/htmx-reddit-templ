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

type Pages struct {
	Home     httprouter.Handle
	Post     httprouter.Handle
	AllPosts httprouter.Handle
	NewPost  httprouter.Handle
	NewUser  httprouter.Handle
}

func New(
	c service.Comment,
	p service.Post,
	u service.User,
) *Pages {
	return &Pages{
		AllPosts: allPosts(p),
		Post:     postPage(p, c),
		NewPost:  newPost(),
		NewUser:  newUser(),
	}
}

func allPosts(ps service.Post) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		postList, err := ps.GetAll()
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}

		templ.AllPosts(postList).Render(context.TODO(), w)
	}
}

func newPost() httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		templ.NewPost().Render(context.TODO(), w)
	}
}

func newUser() httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var pageData struct {
			ShowPassErr bool
		}
		pageData.ShowPassErr = false
		templ.NewUser().Render(context.TODO(), w)
	}
}

func postPage(postService service.Post, commentService service.Comment) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
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
