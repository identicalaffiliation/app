package rest

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/identicalaffiliation/app/pkg/jwtoken"
)

func authMiddleware(tokenValidator jwtoken.TokenValidator) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, errors.New("auth header required").Error(), http.StatusUnauthorized)

				return
			}

			parts := strings.Split(authHeader, " ")
			if len(parts) != 2 || parts[0] != "Bearer" {
				http.Error(w, errors.New("invalid auth format").Error(), http.StatusUnauthorized)

				return
			}

			token := parts[1]
			claims, err := tokenValidator.ValidateTokenWithClaims(token)
			if err != nil {
				if strings.Contains(err.Error(), "signing method") {
					http.Error(w, errors.New("invalid token signing method").Error(), http.StatusUnauthorized)

					return
				}

				http.Error(w, err.Error(), http.StatusUnauthorized)
				return
			}

			if err := tokenValidator.ValidateClaims(claims); err != nil {
				http.Error(w, err.Error(), http.StatusUnauthorized)

				return
			}

			ctx := context.WithValue(r.Context(), "userID", claims["userID"])
			ctx = context.WithValue(ctx, "email", claims["email"])
			r = r.WithContext(ctx)

			next.ServeHTTP(w, r)
		})
	}
}
