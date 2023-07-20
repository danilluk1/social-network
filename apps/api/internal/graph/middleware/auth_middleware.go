package middleware

import (
	"net/http"

	"github.com/danilluk1/social-network/libs/grpc/generated/auth"
)

func Auth(authClient auth.AuthClient) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c, err := r.Cookie("access-token")
			if err != nil || c == nil {
				next.ServeHTTP(w, r)
				return
			}

			res, err := authClient.ValidateUser(r.Context(), &auth.ValidateUserRequest{
				AccessToken: c.Value,
			})
			if err != nil {

			}
		})
	}
}
