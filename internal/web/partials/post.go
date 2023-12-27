package components

import (
	"errors"
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/templ"
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
			req.Context(),
			req.FormValue("title"),
			req.FormValue("body"),
		)
		if err != nil {
			log.Error("Error creating post", "error", err)
			if errors.Is(err, service.ErrSqlError) {
				templ.ErrorSql().Render(req.Context(), w)
			} else {
				templ.ErrorUnknownError().Render(req.Context(), w)
			}

			w.WriteHeader(http.StatusOK)
			return
		}
		log.Info("Created Post", "title", req.FormValue("title"))

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

		if err = svc.Delete(req.Context(), int(id)); err != nil {
			log.Error("could't delete post", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Redirect", "/posts")
		w.WriteHeader(http.StatusOK)
	}
}
