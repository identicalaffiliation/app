package rest

import (
	"context"
	"net/http"

	"github.com/identicalaffiliation/app/internal/config"
)

type HTTPController interface {
	Serve() error
	Shutdown(ctx context.Context) error
}

type httpServer struct {
	server *http.Server
}

func NewHTTPServer(r *Router, cfg *config.AppConfig) HTTPController {
	server := http.Server{
		Addr:    cfg.HTTPServer.Port,
		Handler: r.mux,
	}

	return &httpServer{server: &server}
}

func (http *httpServer) Serve() error {
	return http.server.ListenAndServe()
}

func (http *httpServer) Shutdown(ctx context.Context) error {
	return http.server.Shutdown(ctx)
}
