package post

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/models/post"
	"htmx-reddit/internal/render"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Add    httprouter.Handle
	Delete httprouter.Handle
}

func New(posts post.Model, renderer render.Renderer) *Handler {
	return &Handler{
		Add:    addEndpoint(posts, renderer),
		Delete: deleteEndpoint(posts, renderer),
	}
}

func addEndpoint(model post.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		_, err := model.Add(Data{
			Title: req.FormValue("title"),
			Body:  req.FormValue("body"),
		}.asDBData())
		if err != nil {
			log.Error("getting post id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.FormValue("redirect") == "true" {
			w.Header().Set("HX-Redirect", "/posts#")
		}
		w.WriteHeader(http.StatusOK)
	}
}

func deleteEndpoint(posts post.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = posts.Delete(int(id)); err != nil {
			log.Error("could't delete post", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Redirect", "/posts")
		w.WriteHeader(http.StatusOK)
	}
}
