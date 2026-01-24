package network

import (
	"net/http"
)

type NetworkWriter interface {
	ErrorResponse(w http.ResponseWriter, err error, code int)
	CreatedResponse(w http.ResponseWriter)
	UserFoundResponse(w http.ResponseWriter, userData []byte)
	Response(w http.ResponseWriter)
	AuthResponse(w http.ResponseWriter, authData []byte)
}

type networkWriter struct{}

func NewNetworkWriter() NetworkWriter { return &networkWriter{} }

func (nw *networkWriter) ErrorResponse(w http.ResponseWriter, err error, code int) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "application/json")

	response := []byte(err.Error())
	w.Write(response)
}

func (nw *networkWriter) CreatedResponse(w http.ResponseWriter) {
	w.WriteHeader(http.StatusCreated)
}

func (nw *networkWriter) UserFoundResponse(w http.ResponseWriter, userData []byte) {
	w.WriteHeader(http.StatusFound)
	w.Header().Set("Content-Type", "application/json")
	w.Write(userData)
}

func (nw *networkWriter) Response(w http.ResponseWriter) {
	w.WriteHeader(http.StatusOK)
}

func (nw *networkWriter) AuthResponse(w http.ResponseWriter, authData []byte) {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	w.Write(authData)
}
