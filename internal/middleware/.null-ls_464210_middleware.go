package middleware

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs"
)

type Middleware func(http.Handler) http.Handler

func Group(mws ...Middleware) Middleware {
	return func(next http.Handler) http.Handler {
		for i := 0; i < len(mws); i++ {
			next = mws[i](next)
		}
		return next
	}
}

func authenticateMiddleware(sess *scs.Session) Middleware {

	return func(next http.Handler) Middleware {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				id := app.sessionManager.GetInt(r.Context(), "authenticatedUserID")
				if id == 0 {
					next.ServeHTTP(w, r)
					return
				}

				exists, err := app.users.Exists(id)
				if err != nil {
					app.serverError(w, err)
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
}
