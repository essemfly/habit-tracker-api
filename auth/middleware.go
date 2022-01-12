package auth

import (
	"context"
	"net/http"

	"github.com/lessbutter/habit-tracker-api/repository"
)

var userCtxKey = &contextKey{"user"}

type contextKey struct {
	name string
}

func Middleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			header := r.Header.Get("Authorization")
			// Allow unauthenticated users in
			if header == "" {
				next.ServeHTTP(w, r)
				return
			}
			//validate jwt token
			tokenStr := header
			email, err := ParseToken(tokenStr)
			// Expired된 경우 여기로 통과해서 들어옴
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			// create user and check if user exists in db
			user, err := repository.GetUserByEmail(email)
			if err != nil {
				next.ServeHTTP(w, r)
				return
			}

			// put it in context
			ctx := context.WithValue(r.Context(), userCtxKey, user)

			// and call the next with our new context
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

// ForContext finds the user from the context. REQUIRES Middleware to have run.
func ForContext(ctx context.Context) *repository.UserDAO {
	raw, _ := ctx.Value(userCtxKey).(*repository.UserDAO)

	return raw
}
