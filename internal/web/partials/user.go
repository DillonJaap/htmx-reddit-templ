package components

import (
	"htmx-reddit/internal/convert"
	"htmx-reddit/internal/service"
	"htmx-reddit/internal/templ"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/julienschmidt/httprouter"
)

type user struct {
	Add       http.HandlerFunc
	Delete    http.HandlerFunc
	CheckPass http.HandlerFunc
	Login     http.HandlerFunc
	Logout    http.HandlerFunc
}

func newUser(users service.User) *user {
	return &user{
		Add:       addUser(users),
		Delete:    deleteUser(users),
		CheckPass: checkPassword(users),
		Login:     login(users),
		Logout:    logout(users),
	}
}

func checkPassword(model service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			templ.PasswordsDoNotMatch().Render(req.Context(), w)
			return
		}

		w.WriteHeader(http.StatusNoContent)
	}
}

func addUser(model service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// confirm password
		pass := req.FormValue("password")
		confirmPass := req.FormValue("password-confirm")
		if pass != confirmPass {
			templ.PasswordsDoNotMatch().Render(req.Context(), w)
			return
		}

		err := model.Add(req.FormValue("username"), pass)
		if err != nil {
			log.Error("getting user id", "error", err)
			templ.FailedToCreateUser().Render(req.Context(), w)
			return
		}

		w.Header().Set("HX-Redirect", "/users/login")
		w.WriteHeader(http.StatusOK)
	}
}

func deleteUser(user service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		p := httprouter.ParamsFromContext(req.Context())

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
		w.WriteHeader(http.StatusSeeOther)
	}
}

func login(user service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		// confirm password
		pass := req.FormValue("password")
		username := req.FormValue("username")

		err := user.Login(req.Context(), username, pass)
		if err != nil {
			templ.FailedToLogin().Render(req.Context(), w)
			return
		}

		w.Header().Set("HX-Redirect", "/posts#")
		w.WriteHeader(http.StatusSeeOther)
	}
}

func logout(user service.User) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		user.Logout(req.Context())
		w.Header().Set("HX-Refresh", "true")
		w.WriteHeader(http.StatusSeeOther)
	}
}
