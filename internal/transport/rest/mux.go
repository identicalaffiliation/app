package rest

import (
	"github.com/go-chi/chi/v5"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/pkg/jwtoken"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(cfg *config.AppConfig, ah AuthHandler, uh UserHandler) *Router {
	mux := chi.NewRouter()
	tokenValidator := jwtoken.NewTokenValidator(cfg.JWTSecret)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(func(r chi.Router) {
		r.Use(authMiddleware(tokenValidator))

		r.Get("/api/users/me", uh.MyProfile)
	})

	mux.Group(func(r chi.Router) {
		r.Post("/api/register", ah.SignUp)
		r.Post("/api/login", ah.SignIn)
	})

	return &Router{mux: mux}
}
