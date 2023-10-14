package user

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/models/user"
	"htmx-reddit/internal/render"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Add       httprouter.Handle
	Delete    httprouter.Handle
	CheckPass httprouter.Handle
}

func New(users user.Model, renderer render.Renderer) *Handler {
	return &Handler{
		Add:       addEndpoint(users, renderer),
		Delete:    deleteEndpoint(users, renderer),
		CheckPass: checkPassEndpoint(users, renderer),
	}
}

func checkPassEndpoint(model user.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			r.RenderComponent(w, http.StatusOK, "pass-err", true)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func addEndpoint(model user.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// confirm password
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			r.RenderComponent(w, http.StatusOK, "pass-err", true)
			return
		}

		err := model.Add(UserData{
			Name:     req.FormValue("username"),
			Password: pass,
		}.asDBData())
		if err != nil {
			log.Error("getting user id", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if req.FormValue("redirect") == "true" {
			w.Header().Set("HX-Redirect", "/posts#")
		}
		w.WriteHeader(http.StatusOK)
	}
}

func deleteEndpoint(user user.Model, r render.Renderer) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		id, err := convert.Int(p.ByName("id"))
		if err != nil {
			log.Error("couldn't convert int", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		if err = user.Delete(int(id)); err != nil {
			log.Error("could't delete user", "error", err)
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		w.Header().Set("HX-Redirect", "/posts#")
		w.WriteHeader(http.StatusOK)
	}
}
