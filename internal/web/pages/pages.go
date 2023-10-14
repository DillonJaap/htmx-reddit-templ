package pages

import (
	"htmx-reddit/internal/controller"
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/models/comment"
	"htmx-reddit/internal/models/post"
	"htmx-reddit/internal/models/user"
	"htmx-reddit/internal/render"
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
	r render.Renderer,
) *Pages {
	return &Pages{
		Home:     home(c, r),
		AllPosts: allPosts(p, r),
		Post:     postPage(p, c, r),
		NewPost:  newPost(r),
		NewUser:  newUser(r),
	}
}

func home(model comment.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		commentList := webCmt.GetByParentID(model, 0)
		r.RenderPage(w, http.StatusOK, "index.html", commentList)
	}
}

func allPosts(model post.Model, r render.Renderer) httprouter.Handle {
	getAll := controller.GetAll("post", model.GetAll, webPost.AsViewData)
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		postList, err := getAll()
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}

		r.RenderPage(w, http.StatusOK, "posts.html", postList)
	}
}

func newPost(r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		r.RenderPage(w, http.StatusOK, "post_create.html", "")
	}
}

func newUser(r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var pageData struct {
			ShowPassErr bool
		}
		pageData.ShowPassErr = false
		r.RenderPage(w, http.StatusOK, "user_create.html", pageData)
	}
}

func postPage(postModel post.Model, cmtModel comment.Model, r render.Renderer) httprouter.Handle {
	get := controller.Get("post", postModel.Get, webPost.AsViewData)
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		var pageData struct {
			Post     webPost.Data
			Comments []webCmt.CommentData
		}

		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("failed to convert to int", "error", err)
		}

		pageData.Post, err = get(id)
		if err != nil {
			log.Error("failed to get posts", "error", err)
		}
		log.Printf("data: %+v", pageData)

		pageData.Comments = webCmt.GetByPostID(cmtModel, pageData.Post.ID)

		r.RenderPage(w, http.StatusOK, "post.html", pageData)
	}
}
