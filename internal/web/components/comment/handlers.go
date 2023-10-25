package comment

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
	Add       httprouter.Handle
	Reply     httprouter.Handle
	Delete    httprouter.Handle
	HideReply httprouter.Handle
	ShowReply httprouter.Handle
}

func New(comments service.Comment) *Handler {
	return &Handler{
		Add:       add(comments),
		Reply:     reply(comments),
		Delete:    delete(comments),
		HideReply: hideReplyBox(),
		ShowReply: showReplyBox(),
	}
}

func add(svc service.Comment) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		postID, err := convert.Int(req.FormValue("post-id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := svc.Add(
			postID,
			req.FormValue("comment"),
		)
		if err != nil {
			log.Error("getting comment id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		comment, err := svc.Get(id)
		if err != nil {
			log.Error("quering comment", "error", err, "id", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		templ.Comment(comment).Render(context.TODO(), w)
	}
}

func delete(svc service.Comment) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = svc.Delete(int(id)); err != nil {
			log.Error("could't delete comment", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func reply(svc service.Comment) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		parent_id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("no id received")
			return
		}

		id, err := svc.Add(
			parent_id,
			req.FormValue("comment"),
		)
		if err != nil {
			log.Error("getting comment id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		comment, err := svc.Get(id)
		if err != nil {
			log.Error("quering comment", "error", err, "id", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("HX-Trigger", "hide")
		templ.Comment(comment).Render(context.TODO(), w)
	}
}

func hideReplyBox() httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		strID := req.URL.Query().Get("id")
		id, _ := convert.Int(strID)
		templ.ReplyBox(id, false).Render(context.TODO(), w)
	}
}

func showReplyBox() httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		strID := req.URL.Query().Get("id")
		id, _ := convert.Int(strID)
		templ.ReplyBox(id, true).Render(context.TODO(), w)
	}
}
