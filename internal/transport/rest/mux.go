package rest

import (
	"github.com/go-chi/chi"
	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/pkg/jwtoken"
)

type Router struct {
	mux *chi.Mux
}

func NewRouter(cfg *config.AppConfig) *Router {
	mux := chi.NewRouter()
	tokenValidator := jwtoken.NewTokenValidator(cfg.JWTSecret)

	mux.Use(authMiddleware(tokenValidator))
	return &Router{mux: mux}
}
