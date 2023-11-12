package middleware

import (
	"context"
	"htmx-reddit/internal/helpers"
	"htmx-reddit/internal/service"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Middleware func(http.Handler) http.Handler

func Join(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := len(mws) - 1; i >= 0; i-- {
			next = mws[i](next)
		}
		return next
	}
}

func Authenticate(sess *scs.Session, user service.User) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := sess.GetInt(r.Context(), "authenticatedUserID")
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}
			exists, err := user.Exists(id)
			if err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}

			if exists {
				ctx := context.WithValue(r.Context(), "isAuthenticated", true)
				r = r.WithContext(ctx)
			}

			next.ServeHTTP(w, r)
		})
	}
}

func RequireAuthentication() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if !helpers.IsAuthenticated(r) {
				w.Header().Set("HX-Redirect", "/users/login")
				return
			}
			w.Header().Add("Cache-Control", "no-store")
			next.ServeHTTP(w, r)
		})
	}
}
