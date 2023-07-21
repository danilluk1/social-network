package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
)

// A private key for context that only this package can access. This is important
// to prevent collisions between different context uses
var userCtxKey = &contextKey{name: "user"}

type contextKey struct {
	name string
}

func Auth(authClient auth.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := TokenFromHTTPRequestgo(r)
			if c == "" {
				next.ServeHTTP(w, r)
				return
			}

			res, err := authClient.ValidateUser(r.Context(), &auth.ValidateUserRequest{
				AccessToken: c,
			})
			if err != nil {
				http.Error(w, "Invalid token", http.StatusForbidden)
				return
			}

			ctx := context.WithValue(r.Context(), userCtxKey, res.Username)
			r = r.WithContext(ctx)
			next.ServeHTTP(w, r)
		})
	}
}

func TokenFromHTTPRequestgo(r *http.Request) string {
	reqToken := r.Header.Get("Authorization")
	var tokenString string

	splitToken := strings.Split(reqToken, "Bearer ")
	if len(splitToken) > 1 {
		tokenString = splitToken[1]
	}
	return tokenString
}

func ForContext(ctx context.Context) *string {
	raw := ctx.Value(userCtxKey).(string)
	return &raw
}
