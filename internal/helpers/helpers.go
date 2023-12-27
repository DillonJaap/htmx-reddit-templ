package helpers

import (
	"context"
	"net/http"

	"github.com/alexedwards/scs/v2"
)

func IsAuthenticated(r *http.Request) bool {
	isAuthenticated, ok := r.Context().Value("isAuthenticated").(bool)
	if !ok {
		return false
	}
	return isAuthenticated
}

func IsLoggedInUser(ctx context.Context, sess *scs.SessionManager, wantID int) bool {
	return sess.GetInt(ctx, "authenticatedUserID") == wantID
}

func IsLoggedIn(ctx context.Context, sess *scs.SessionManager) bool {
	return sess.GetInt(ctx, "authenticatedUserID") != 0
}
