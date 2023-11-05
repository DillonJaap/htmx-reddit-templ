package components

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/service"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type post struct {
	Add    http.HandlerFunc
	Delete http.HandlerFunc
}

func newPost(posts service.Post) *post {
	return &post{
		Add:    addPost(posts),
		Delete: deletePost(posts),
	}
}

func addPost(svc service.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		err := svc.Add(
			req.FormValue("title"),
			req.FormValue("body"),
		)
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

func deletePost(svc service.Post) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p := httprouter.ParamsFromContext(req.Context())
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = svc.Delete(int(id)); err != nil {
			log.Error("could't delete post", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Redirect", "/posts")
		w.WriteHeader(http.StatusOK)
	}
}
