package rest

import (
	"net/http"

	se "github.com/identicalaffiliation/app/internal/service/entity"
)

type AuthHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
}

type authHandler struct {
	authService se.AuthUseCases
}

func NewAuthHandler(as se.AuthUseCases) AuthHandler {
	return &authHandler{
		authService: as,
	}
}

func (ah *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	//TODO
}

func (ah *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	//TODO
}
