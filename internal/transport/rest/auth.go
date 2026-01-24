package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/identicalaffiliation/app/internal/dto"
	se "github.com/identicalaffiliation/app/internal/service/entity"
	"github.com/identicalaffiliation/app/pkg/network"
)

type authHandler struct {
	authService se.AuthUseCases
	nw          network.NetworkWriter
}

func NewAuthHandler(as se.AuthUseCases) AuthHandler {
	nw := network.NewNetworkWriter()

	return &authHandler{
		authService: as,
		nw:          nw,
	}
}

func (ah *authHandler) SignUp(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ah.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.UserRegisterRequest
	if err := json.Unmarshal(body, &request); err != nil {
		ah.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	if err := ah.authService.Register(r.Context(), &request); err != nil {
		ah.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	ah.nw.CreatedResponse(w)
}

func (ah *authHandler) SignIn(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		ah.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		ah.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.UserLoginRequest
	if err := json.Unmarshal(body, &request); err != nil {
		ah.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	response, err := ah.authService.Login(r.Context(), &request)
	if err != nil {
		ah.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	authData, err := json.Marshal(response)
	if err != nil {
		ah.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	ah.nw.AuthResponse(w, authData)
}
