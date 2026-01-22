package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/dto"
	re "github.com/identicalaffiliation/app/internal/repository/entity"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	se "github.com/identicalaffiliation/app/internal/service/entity"
	"github.com/identicalaffiliation/app/pkg/hash"
)

type userService struct {
	userRepo  psql.UserRepository
	validator *se.Validator
	hasher    hash.Hasher
}

func NewUserService(ur psql.UserRepository) se.UserUseCases {
	v := se.InitValidator()
	h := hash.NewHasher()

	return &userService{
		userRepo:  ur,
		validator: v,
		hasher:    h,
	}
}

func (us *userService) CreateUser(ctx context.Context, userRequest *dto.UserCreateRequest) error {
	if err := us.validator.UserCreateValidate(userRequest); err != nil {
		return fmt.Errorf(se.ErrInvalidCreateUserRequest.Error(), err)
	}

	uuid := uuid.New()

	hashedPassword, err := us.hasher.HashPassword(userRequest.Password)
	if err != nil {
		return err
	}

	user := &re.User{
		ID:       uuid,
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: hashedPassword,
	}

	return us.userRepo.Create(ctx, user)
}
