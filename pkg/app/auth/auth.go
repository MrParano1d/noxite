package auth

import (
	"context"
	"log"
	"net/http"
	"strings"

	"github.com/mrparano1d/noxite/pkg/core"
	"github.com/mrparano1d/noxite/pkg/core/entities"
	"github.com/mrparano1d/noxite/pkg/core/services"
)

type AuthContextKey string

const (
	AuthContextSessionKey AuthContextKey = "session"
	AuthContextUserKey    AuthContextKey = "user"
)

func AuthMiddleware(coreApp *core.ApplicationCore) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			ctx := req.Context()
			token := req.Header.Get("Authorization")

			token = strings.Replace(token, "Bearer ", "", 1)

			if token == "" {
				log.Println("no authorization token provided")
				http.Error(w, "no authorization token provided", http.StatusUnauthorized)
				return
			}

			err := coreApp.SessionService().ValidateToken(ctx, token)
			if err != nil {
				log.Printf("invalid token: %s\n", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			user, err := services.SessionValueFromService[entities.User](coreApp.SessionService(), ctx, token, "user")
			if err != nil {
				log.Printf("failed to get user from session: %s\n", err)
				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			ctx = context.WithValue(ctx, AuthContextSessionKey, token)
			ctx = context.WithValue(ctx, AuthContextUserKey, user)

			next.ServeHTTP(w, req.WithContext(ctx))
		})
	}
}

func GetSessionFromContext(ctx context.Context) string {
	return ctx.Value(AuthContextSessionKey).(string)
}

func GetUserFromContext(ctx context.Context) *entities.User {
	return ctx.Value(AuthContextUserKey).(*entities.User)
}
