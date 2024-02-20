package graphql

import (
	"context"
	"net/http"
	"strings"
)

type AuthContextKey string

const (
	AuthContextKeyToken AuthContextKey = "token"
)

func TokenMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bearer := r.Header.Get("Authorization")

		token := strings.TrimPrefix(bearer, "Bearer ")

		ctx := r.Context()

		if token != "" {
			ctx = context.WithValue(r.Context(), AuthContextKeyToken, token)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func TokenFromContext(ctx context.Context) (string, bool) {
	token, ok := ctx.Value(AuthContextKeyToken).(string)
	if !ok {
		return "", false
	}
	return token, true
}

/*func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		role := RolePublic

	})
}*/
