package pages

import (
	"context"
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/models/comment"
	"htmx-reddit/internal/models/post"
	"htmx-reddit/internal/models/user"
	"htmx-reddit/internal/templ"
	webCmt "htmx-reddit/internal/web/components/comment"
	webPost "htmx-reddit/internal/web/components/post"
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
	c comment.Model,
	p post.Model,
	u user.Model,
) *Pages {
	return &Pages{
		AllPosts: allPosts(p),
		Post:     postPage(p, c),
		NewPost:  newPost(),
		NewUser:  newUser(),
	}
}

func allPosts(model post.Model) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		postList, err := webPost.GetAll(model)
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

func postPage(postModel post.Model, cmtModel comment.Model) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("failed to convert to int", "error", err)
		}

		posts, err := webPost.Get(postModel, id)
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}

		comments := webCmt.GetByPostID(cmtModel, posts.ID)
		templ.Post(posts, comments).Render(context.TODO(), w)
	}
}
