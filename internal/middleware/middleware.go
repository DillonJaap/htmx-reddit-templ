package middleware

import (
	"context"
	"htmx-reddit/internal/service"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

type Middleware func(http.Handler) http.Handler

func Join(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(mws); i++ {
			next = mws[i](next)
		}
		return next
	}
}

func AuthenticateMiddleware(sess *scs.Session, user service.User) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			id := sess.GetInt(r.Context(), "authenticatedUserID")
			if id == 0 {
				next.ServeHTTP(w, r)
				return
			}

			exists, err := user.Exists(id)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
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
