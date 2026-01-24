package rest

import (
	"net/http"
)

type AuthHandler interface {
	SignUp(w http.ResponseWriter, r *http.Request)
	SignIn(w http.ResponseWriter, r *http.Request)
}

type UserHandler interface {
	MyProfile(w http.ResponseWriter, r *http.Request)
}
