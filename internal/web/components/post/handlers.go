package post

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/service"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Add    httprouter.Handle
	Delete httprouter.Handle
}

func New(posts service.Post) *Handler {
	return &Handler{
		Add:    add(posts),
		Delete: delete(posts),
	}
}

func add(svc service.Post) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
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

func delete(svc service.Post) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
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
