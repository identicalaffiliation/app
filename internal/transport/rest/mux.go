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

func NewRouter(cfg *config.AppConfig, ah AuthHandler, uh UserHandler, th TodoHandler) *Router {
	mux := chi.NewRouter()
	tokenValidator := jwtoken.NewTokenValidator(cfg.JWTSecret)

	mux.Use(middleware.Logger)
	mux.Use(middleware.Recoverer)

	mux.Group(func(r chi.Router) {
		r.Post("/api/register", ah.SignUp)
		r.Post("/api/login", ah.SignIn)
	})

	mux.Group(func(r chi.Router) {
		r.Use(authMiddleware(tokenValidator))

		r.Route("/api/users", func(r chi.Router) {

			r.Route("/me", func(r chi.Router) {
				r.Get("/", uh.MyProfile)
				r.Patch("/name", uh.ChangeMyName)
				r.Patch("/email", uh.ChangeMyEmail)
				r.Patch("/password", uh.ChangeMyPassword)

				r.Route("/todos", func(r chi.Router) {
					r.Post("/", th.NewTodo)
					r.Get("/", th.MyTodos)

					r.Route("/{todoID}", func(r chi.Router) {
						r.Get("/", th.MyTodo)
						r.Patch("/content", th.ChangeTodoContent)
						r.Patch("/status", th.ChangeTodoStatus)
						r.Delete("/", th.DeleteTodo)
					})
				})
			})
		})
	})

	return &Router{mux: mux}
}
