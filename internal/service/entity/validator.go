package entity

import (
	"github.com/go-playground/validator"
	"github.com/identicalaffiliation/app/internal/dto"
)

type Validator struct {
	Validator *validator.Validate
}

func InitValidator() *Validator {
	v := validator.New()

	return &Validator{
		Validator: v,
	}
}

func (v *Validator) UserRegisterRequestValidate(user *dto.UserRegisterRequest) error {
	return v.Validator.Struct(user)
}

func (v *Validator) UserLoginRequestValidate(user *dto.UserLoginRequest) error {
	return v.Validator.Struct(user)
}

func (v *Validator) UserChangeNameReguestValidate(userChangeReguest *dto.ChangeUserNameRequest) error {
	return v.Validator.Struct(userChangeReguest)
}

func (v *Validator) UserChangeEmailRequestValidate(userChangeReguest *dto.ChangeUserEmailRequest) error {
	return v.Validator.Struct(userChangeReguest)
}

func (v *Validator) UserChangePasswordRequestValidate(userChangeReguest *dto.ChangeUserPasswordRequest) error {
	return v.Validator.Struct(userChangeReguest)
}

func (v *Validator) TodoCreateRequestValidate(todoRequest *dto.TodoCreateRequest) error {
	if todoRequest.Status == "done" {
		return ErrInvalidTodoStatus
	}

	if todoRequest.Status == "" {
		todoRequest.Status = "todo"
	}

	return v.Validator.Struct(todoRequest)
}

func (v *Validator) TodoContentChangeRequest(todoChangeRequest *dto.TodoContentChangeRequest) error {
	return v.Validator.Struct(todoChangeRequest)
}

func (v *Validator) TodoStatusChangeRequest(todoChangeRequest *dto.TodoStatusChangeRequest) error {
	return v.Validator.Struct(todoChangeRequest)
}
