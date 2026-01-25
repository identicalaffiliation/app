package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/identicalaffiliation/app/internal/config"
	"github.com/identicalaffiliation/app/internal/logger"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	"github.com/identicalaffiliation/app/internal/service"
	"github.com/identicalaffiliation/app/internal/transport/rest"
	"github.com/identicalaffiliation/app/pkg/parse"
)

func main() {
	path := parse.FlagInit()

	cfg := config.MustLoadConfig(path)

	logger := logger.NewLogger()

	db := psql.NewPostgres()
	defer db.Close()

	db.MustInit(cfg)

	userRepo := psql.NewUserRepository(db, logger)
	todoRepo := psql.NewTodoRepository(db, logger)
	userSerivce := service.NewUserService(userRepo)
	todoService := service.NewTodoService(userRepo, todoRepo)
	authService := service.NewAuthService(userRepo, cfg.JWTSecret)
	authHandler := rest.NewAuthHandler(authService)
	userHandler := rest.NewUserHandler(userSerivce)
	todoHandler := rest.NewTodoHandler(todoService)

	r := rest.NewRouter(cfg, authHandler, userHandler, todoHandler)
	s := rest.NewHTTPServer(r, cfg)

	go func() {

		log.Println("server started")
		if err := s.Serve(); err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := s.Shutdown(ctx); err != nil {
		log.Fatal(err)
	}

	log.Println("server stopped gracefully")
}
