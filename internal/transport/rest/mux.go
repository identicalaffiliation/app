package rest

import "github.com/go-chi/chi"

type Router struct {
	mux *chi.Mux
}

func NewRouter() *Router {
	mux := chi.NewRouter()

	return &Router{mux: mux}
}
