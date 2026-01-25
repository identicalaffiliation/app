package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/dto"
	re "github.com/identicalaffiliation/app/internal/repository/entity"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	se "github.com/identicalaffiliation/app/internal/service/entity"
)

type todoService struct {
	userRepo  psql.UserRepository
	todoRepo  psql.TodoRepository
	validator *se.Validator
}

func NewTodoService(ur psql.UserRepository, tr psql.TodoRepository) se.TodoUseCases {
	v := se.InitValidator()

	return &todoService{
		userRepo:  ur,
		todoRepo:  tr,
		validator: v,
	}
}

func (ts *todoService) CreateTodo(ctx context.Context, todoRequest *dto.TodoCreateRequest) error {
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if userID == uuid.Nil {
		return se.ErrInvalidUserID
	}

	if err := ts.validator.TodoCreateRequestValidate(todoRequest); err != nil {
		return fmt.Errorf("todo validate: %w", err)
	}

	todoID := uuid.New()

	todo := &re.Todo{
		ID:      todoID,
		UserID:  userID,
		Content: todoRequest.Content,
		Status:  todoRequest.Status,
	}

	return ts.todoRepo.Create(ctx, todo)
}

func (ts *todoService) GetTodo(ctx context.Context, todoID uuid.UUID) (*dto.TodoResponse, error) {
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if userID == uuid.Nil {
		return nil, se.ErrInvalidUserID
	}

	if todoID == uuid.Nil {
		return nil, se.ErrInvalidTodoID
	}

	todo, err := ts.todoRepo.GetTodoByUserID(ctx, todoID, userID)
	if err != nil {
		return nil, err
	}

	return &dto.TodoResponse{
		Content:   todo.Content,
		Status:    todo.Status,
		CreatedAt: todo.CreatedAt,
		UpdatedAt: todo.UpdatedAt,
	}, nil
}

func (ts *todoService) GetTodos(ctx context.Context) ([]*dto.TodoResponse, error) {
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if userID == uuid.Nil {
		return nil, se.ErrInvalidUserID
	}

	todos, err := ts.todoRepo.GetTodosByUserID(ctx, userID)
	if err != nil {
		return nil, err
	}

	response := ts.todosToResponse(todos)

	return response, nil
}

func (ts *todoService) todosToResponse(todos []*re.Todo) []*dto.TodoResponse {
	respone := make([]*dto.TodoResponse, 0, len(todos))
	for _, todo := range todos {
		respone = append(respone, &dto.TodoResponse{
			Content:   todo.Content,
			Status:    todo.Status,
			CreatedAt: todo.CreatedAt,
			UpdatedAt: todo.UpdatedAt,
		})
	}

	return respone
}

func (ts *todoService) ChangeContent(ctx context.Context, changeContentRequest *dto.TodoContentChangeRequest) error {
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if userID == uuid.Nil {
		return se.ErrInvalidUserID
	}

	if changeContentRequest.TodoID == uuid.Nil {
		return se.ErrInvalidTodoID
	}

	if err := ts.validator.TodoContentChangeRequest(changeContentRequest); err != nil {
		return err
	}

	return ts.todoRepo.UpdateContent(ctx, changeContentRequest.NewContent, changeContentRequest.TodoID, userID)
}

func (ts *todoService) ChangeStatus(ctx context.Context, changeStatusRequest *dto.TodoStatusChangeRequest) error {
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if userID == uuid.Nil {
		return se.ErrInvalidUserID
	}

	if changeStatusRequest.TodoID == uuid.Nil {
		return se.ErrInvalidTodoID
	}

	if err := ts.validator.TodoStatusChangeRequest(changeStatusRequest); err != nil {
		return err
	}

	return ts.todoRepo.UpdateStatus(ctx, psql.TodoStatus(changeStatusRequest.NewStatus), changeStatusRequest.TodoID, userID)
}

func (ts *todoService) DeleteTodo(ctx context.Context, todoID uuid.UUID) error {
	userID, _ := uuid.Parse(ctx.Value("userID").(string))
	if userID == uuid.Nil {
		return se.ErrInvalidUserID
	}

	if todoID == uuid.Nil {
		return se.ErrInvalidUserID
	}

	return ts.todoRepo.Delete(ctx, todoID, userID)
}
