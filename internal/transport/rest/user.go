package rest

import (
	"encoding/json"
	"net/http"

	"github.com/google/uuid"
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

	userIDstr := r.Context().Value("userID").(string)
	userID, _ := uuid.Parse(userIDstr)

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
