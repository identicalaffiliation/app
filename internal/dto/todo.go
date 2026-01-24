package dto

import (
	"time"

	"github.com/google/uuid"
)

type (
	TodoCreateRequest struct {
		Content string `json:"content" validate:"required"`
		Status  string `json:"status" validate:"required,oneof=todo process"`
	}

	TodoResponse struct {
		Content   string    `json:"content"`
		Status    string    `json:"status"`
		CreatedAt time.Time `json:"createdAt"`
		UpdatedAt time.Time `json:"updatedAt"`
	}

	TodoContentChangeRequest struct {
		TodoID     uuid.UUID `json:"todoID" validate:"required"`
		NewContent string    `json:"content" validate:"required"`
	}

	TodoStatusChangeRequest struct {
		TodoID    uuid.UUID `json:"todoID" validate:"required"`
		NewStatus string    `json:"status" validate:"required"`
	}
)
