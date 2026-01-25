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

type todoHandler struct {
	todoService se.TodoUseCases
	nw          network.NetworkWriter
}

func NewTodoHandler(ts se.TodoUseCases) TodoHandler {
	nw := network.NewNetworkWriter()

	return &todoHandler{
		todoService: ts,
		nw:          nw,
	}
}

func (th *todoHandler) NewTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		th.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.TodoCreateRequest
	if err := json.Unmarshal(body, &request); err != nil {
		th.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	if err := th.todoService.CreateTodo(r.Context(), &request); err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	th.nw.CreatedResponse(w)
}

func (th *todoHandler) MyTodo(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		th.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	todoID, err := uuid.Parse(r.PathValue("todoID"))
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	response, err := th.todoService.GetTodo(r.Context(), todoID)
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	todoData, err := json.Marshal(response)
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	th.nw.TodoFoundResponse(w, todoData)
}

func (th *todoHandler) MyTodos(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		th.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	response, err := th.todoService.GetTodos(r.Context())
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	todoData, err := json.Marshal(response)
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	th.nw.TodoFoundResponse(w, todoData)
}

func (th *todoHandler) ChangeTodoContent(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		th.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	todoID, err := uuid.Parse(r.PathValue("todoID"))
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.TodoContentChangeRequest
	if err := json.Unmarshal(body, &request); err != nil {
		th.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	request.TodoID = todoID
	if err := th.todoService.ChangeContent(r.Context(), &request); err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	th.nw.Response(w)
}

func (th *todoHandler) ChangeTodoStatus(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPatch {
		th.nw.ErrorResponse(w, ErrInvalidMethod, http.StatusMethodNotAllowed)

		return
	}

	todoID, err := uuid.Parse(r.PathValue("todoID"))
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		th.nw.ErrorResponse(w, err, http.StatusInternalServerError)

		return
	}

	var request dto.TodoStatusChangeRequest
	if err := json.Unmarshal(body, &request); err != nil {
		th.nw.ErrorResponse(w, ErrInvalidJSONBody, http.StatusBadRequest)

		return
	}

	request.TodoID = todoID
	if err := th.todoService.ChangeStatus(r.Context(), &request); err != nil {
		th.nw.ErrorResponse(w, err, http.StatusBadRequest)

		return
	}

	th.nw.Response(w)
}
