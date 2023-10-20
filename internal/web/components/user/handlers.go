package user

import (
	"context"
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/models/user"
	"htmx-reddit/internal/templ"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type Handler struct {
	Add       httprouter.Handle
	Delete    httprouter.Handle
	CheckPass httprouter.Handle
}

func New(users user.Model) *Handler {
	return &Handler{
		Add:       addEndpoint(users),
		Delete:    deleteEndpoint(users),
		CheckPass: checkPassEndpoint(users),
	}
}

func checkPassEndpoint(model user.Model) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			templ.PasswordsDoNotMatch(true).Render(context.TODO(), w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func addEndpoint(model user.Model) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// confirm password
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			templ.PasswordsDoNotMatch(true).Render(context.TODO(), w)
			return
		}

		err := model.Add(user.User{
			Name:     req.FormValue("username"),
			Password: pass,
		})
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

func deleteEndpoint(user user.Model) httprouter.Handle {
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
