package psql

import (
	"context"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/repository/entity"
)

type TodoRepository interface {
	Create(ctx context.Context, todo *entity.Todo) error
	GetNotesByUserID(ctx context.Context, userID uuid.UUID) ([]*entity.Todo, error)
	GetNoteByUserID(ctx context.Context, todoID, userID uuid.UUID) (*entity.Todo, error)
	UpdateStatus(ctx context.Context, newStatus NoteStatus, todoID, userID uuid.UUID) error
	UpdateContent(ctx context.Context, newContent string, todoID, userID uuid.UUID) error
	Delete(ctx context.Context, todoID, userID uuid.UUID) error
}
