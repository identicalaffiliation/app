package rest

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/dto"
	se "github.com/identicalaffiliation/app/internal/service/entity"
	"github.com/identicalaffiliation/app/pkg/network"
)

type userHandler struct {
	userService se.UserUseCases
	nw          network.NetworkWriter
}

func NewUserHandler(us se.UserUseCases) UserHandler {
	nw := network.NewNetworkWriter()

	return &userHandler{
		userService: us,
		nw:          nw,
	}
}

func (uh *userHandler) MyProfile(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		uh.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	userID, err := uuid.Parse(r.Context().Value("userID").(string))
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	response, err := uh.userService.GetUser(r.Context(), userID)
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	userData, err := json.Marshal(response)
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	uh.nw.UserFoundResponse(w, userData)
}

func (uh *userHandler) ChangeMyName(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		uh.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	userID, err := uuid.Parse(r.Context().Value("userID").(string))
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.ChangeUserNameRequest
	if err := json.Unmarshal(body, &request); err != nil {
		uh.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	request.ID = userID
	if err := uh.userService.ChangeName(r.Context(), &request); err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	uh.nw.Response(w)
}

func (uh *userHandler) ChangeMyEmail(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		uh.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	userID, err := uuid.Parse(r.Context().Value("userID").(string))
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.ChangeUserEmailRequest
	if err := json.Unmarshal(body, &request); err != nil {
		uh.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	request.ID = userID
	if err := uh.userService.ChangeEmail(r.Context(), &request); err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	uh.nw.Response(w)
}

func (uh *userHandler) ChangeMyPassword(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		uh.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	userID, err := uuid.Parse(r.Context().Value("userID").(string))
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.ChangeUserPasswordRequest
	if err := json.Unmarshal(body, &request); err != nil {
		uh.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	request.ID = userID
	if err := uh.userService.ChangePassword(r.Context(), &request); err != nil {
		uh.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	uh.nw.Response(w)
}
