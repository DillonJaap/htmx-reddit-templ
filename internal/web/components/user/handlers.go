package user

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
	Delete    httprouter.Handle
	CheckPass httprouter.Handle
}

func New(users service.User) *Handler {
	return &Handler{
		Add:       add(users),
		Delete:    delete(users),
		CheckPass: checkPassword(users),
	}
}

func checkPassword(model service.User) httprouter.Handle {
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

func add(model service.User) httprouter.Handle {
	return func(w http.ResponseWriter, req *http.Request, p httprouter.Params) {
		// confirm password
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			templ.PasswordsDoNotMatch(true).Render(context.TODO(), w)
			return
		}

		err := model.Add(req.FormValue("username"), pass)
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

func delete(user service.User) httprouter.Handle {
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
