package components

import (
	"context"
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/templ"
	"net/http"

	"github.com/alexedwards/scs/v2"
	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type comment struct {
	Add       http.HandlerFunc
	Reply     http.HandlerFunc
	Delete    http.HandlerFunc
	HideReply http.HandlerFunc
	ShowReply http.HandlerFunc
}

func newComment(comments service.Comment, sess *scs.SessionManager) *comment {
	return &comment{
		Add:       addComment(comments, sess),
		Reply:     reply(comments, sess),
		Delete:    deleteComment(comments),
		HideReply: hideReplyBox(),
		ShowReply: showReplyBox(),
	}
}

func addComment(svc service.Comment, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		postID, err := convert.Int(req.FormValue("post-id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := svc.Add(
			req.Context(),
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
		templ.Comment(
			comment,
			sess.GetString(req.Context(), "username"),
		).Render(context.TODO(), w)
	}
}

func deleteComment(svc service.Comment) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p := httprouter.ParamsFromContext(req.Context())
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = svc.Delete(req.Context(), int(id)); err != nil {
			log.Error("could't delete comment", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func reply(svc service.Comment, sess *scs.SessionManager) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p := httprouter.ParamsFromContext(req.Context())
		parentID, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("no id received")
			return
		}

		id, err := svc.Reply(
			req.Context(),
			parentID,
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
		templ.Comment(
			comment,
			sess.GetString(req.Context(), "username"),
		).Render(context.TODO(), w)
	}
}

func hideReplyBox() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		strID := req.URL.Query().Get("id")
		id, _ := convert.Int(strID)
		templ.ReplyBox(id, false).Render(context.TODO(), w)
	}
}

func showReplyBox() http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		strID := req.URL.Query().Get("id")
		id, _ := convert.Int(strID)
		templ.ReplyBox(id, true).Render(context.TODO(), w)
	}
}
