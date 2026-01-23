package service

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/dto"
	"github.com/identicalaffiliation/app/internal/repository/entity"
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

func (us *userService) GetUser(ctx context.Context, userID uuid.UUID) (*dto.UserResponse, error) {
	user, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}

	response := &dto.UserResponse{
		ID:    user.ID,
		Email: user.Email,
		Name:  user.Name,
	}

	return response, nil
}

func (us *userService) GetUsers(ctx context.Context) ([]*dto.UserResponse, error) {
	users, err := us.userRepo.GetAllUsers(ctx)
	if err != nil {
		return nil, fmt.Errorf("get users: %w", err)
	}

	response := us.usersToResponse(users)

	return response, nil
}

func (us *userService) usersToResponse(users []*entity.User) []*dto.UserResponse {
	response := make([]*dto.UserResponse, 0, len(users))
	for _, user := range users {
		response = append(response, &dto.UserResponse{
			ID:    user.ID,
			Email: user.Email,
			Name:  user.Name,
		})
	}

	return response
}

func (us *userService) ChangeName(ctx context.Context, changeNameRequest *dto.ChangeUserNameRequest) error {
	if err := us.validator.UserChangeNameReguestValidate(changeNameRequest); err != nil {
		return err
	}

	user, err := us.userRepo.GetByID(ctx, changeNameRequest.ID)
	if err != nil {
		return se.ErrInvalidUserID
	}

	if err := us.hasher.CompareHashAndPassword(user.Password, changeNameRequest.Password); err != nil {
		return se.ErrInvalidPassword
	}

	return us.userRepo.ChangeName(ctx, changeNameRequest.Name, changeNameRequest.ID)
}

func (us *userService) ChangeEmail(ctx context.Context, changeEmailRequest *dto.ChangeUserEmailRequest) error {
	if err := us.validator.UserChangeEmailRequestValidate(changeEmailRequest); err != nil {
		return err
	}

	user, err := us.userRepo.GetByID(ctx, changeEmailRequest.ID)
	if err != nil {
		return se.ErrInvalidUserID
	}

	if err := us.hasher.CompareHashAndPassword(user.Password, changeEmailRequest.Password); err != nil {
		return se.ErrInvalidPassword
	}

	return us.userRepo.ChangeEmail(ctx, changeEmailRequest.Email, changeEmailRequest.ID)
}

func (us *userService) ChangePassword(ctx context.Context, changePasswordRequest *dto.ChangeUserPasswordRequest) error {
	if err := us.validator.UserChangePasswordRequestValidate(changePasswordRequest); err != nil {
		return err
	}

	user, err := us.userRepo.GetByID(ctx, changePasswordRequest.ID)
	if err != nil {
		return se.ErrInvalidUserID
	}

	if err := us.hasher.CompareHashAndPassword(user.Password, changePasswordRequest.OldPassword); err != nil {
		return se.ErrInvalidPassword
	}

	hashedPassword, err := us.hasher.HashPassword(changePasswordRequest.NewPassword)
	if err != nil {
		return err
	}

	return us.userRepo.ChangePassword(ctx, hashedPassword, changePasswordRequest.ID)
}

func (us *userService) DeleteUser(ctx context.Context, userID uuid.UUID) error {
	if userID == uuid.Nil {
		return se.ErrInvalidUserID
	}

	_, err := us.userRepo.GetByID(ctx, userID)
	if err != nil {
		return se.ErrInvalidUserID
	}

	return us.userRepo.Delete(ctx, userID)
}
