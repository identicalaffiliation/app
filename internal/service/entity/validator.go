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
