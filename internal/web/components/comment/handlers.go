package comment

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/models/comment"
	"htmx-reddit/internal/render"
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

func New(comments comment.Model, renderer render.Renderer) *Handler {
	return &Handler{
		Add:       addEndpoint(comments, renderer),
		Reply:     replyEndpoint(comments, renderer),
		Delete:    deleteEndpoint(comments, renderer),
		HideReply: hideReplyBox(renderer),
		ShowReply: showReplyBox(renderer),
	}
}

func addEndpoint(comments comment.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		postID, err := convert.Int(req.FormValue("post-id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		id, err := comments.Add(comment.Comment{
			ParentID:    postID,
			Description: req.FormValue("comment"),
		})
		if err != nil {
			log.Error("getting comment id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		commentItem, err := comments.Get(id)
		if err != nil {
			log.Error("quering comment", "error", err, "id", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		r.RenderComponent(w, http.StatusOK, "comment", asViewData(commentItem))
	}
}

func deleteEndpoint(comments comment.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = comments.Delete(int(id)); err != nil {
			log.Error("could't delete comment", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.WriteHeader(http.StatusOK)
	}
}

func replyEndpoint(comments comment.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		parent_id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("no id received")
			return
		}

		id, err := comments.Add(comment.Comment{
			ParentID:    parent_id,
			Description: req.FormValue("comment"),
		})
		if err != nil {
			log.Error("getting comment id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		commentItem, err := comments.Get(id)
		if err != nil {
			log.Error("quering comment", "error", err, "id", id)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Add("HX-Trigger", "hide")
		r.RenderComponent(w, http.StatusOK, "comment", asViewData(commentItem))
	}
}

func hideReplyBox(r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ID := req.URL.Query().Get("id")
		strID, _ := convert.Int(ID)
		r.RenderComponent(w, http.StatusOK, "reply-box", struct {
			ReplyBoxVisible bool
			ID              int
		}{
			ReplyBoxVisible: false,
			ID:              strID,
		},
		)
	}
}

func showReplyBox(r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		ID := req.URL.Query().Get("id")
		strID, _ := convert.Int(ID)
		r.RenderComponent(w, http.StatusOK, "reply-box", struct {
			ReplyBoxVisible bool
			ID              int
		}{
			ReplyBoxVisible: true,
			ID:              strID,
		},
		)
	}
}
