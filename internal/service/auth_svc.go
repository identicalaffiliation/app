package service

import (
	"context"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/google/uuid"
	"github.com/identicalaffiliation/app/internal/dto"
	re "github.com/identicalaffiliation/app/internal/repository/entity"
	"github.com/identicalaffiliation/app/internal/repository/psql"
	se "github.com/identicalaffiliation/app/internal/service/entity"
	"github.com/identicalaffiliation/app/pkg/hash"
)

type authService struct {
	userRepo  psql.UserRepository
	validator *se.Validator
	hasher    hash.Hasher
	jwtSecret string
}

func NewAuthService(ur psql.UserRepository, secret string) se.AuthUseCases {
	v := se.InitValidator()
	h := hash.NewHasher()

	return &authService{
		userRepo:  ur,
		validator: v,
		hasher:    h,
		jwtSecret: secret,
	}
}

func (as *authService) Register(ctx context.Context, userRequest *dto.UserRegisterRequest) error {
	if err := as.validator.UserRegisterRequestValidate(userRequest); err != nil {
		return err
	}

	uuid := uuid.New()

	hashedPassword, err := as.hasher.HashPassword(userRequest.Password)
	if err != nil {
		return err
	}

	user := &re.User{
		ID:       uuid,
		Name:     userRequest.Name,
		Email:    userRequest.Email,
		Password: hashedPassword,
	}

	return as.userRepo.Create(ctx, user)
}

func (as *authService) Login(ctx context.Context, userRequest *dto.UserLoginRequest) (*dto.AuthResponse, error) {
	if err := as.validator.UserLoginRequestValidate(userRequest); err != nil {
		return nil, err
	}

	user, err := as.userRepo.GetByEmail(ctx, userRequest.Email)
	if err != nil {
		return nil, se.ErrInvalidUserEmail
	}

	if err := as.hasher.CompareHashAndPassword(user.Password, userRequest.Password); err != nil {
		return nil, se.ErrInvalidPassword
	}

	token, expires, err := as.generateToken(user)
	if err != nil {
		return nil, err
	}

	response := &dto.AuthResponse{
		User: &dto.UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
		Token:     token,
		ExpiresAt: expires,
	}

	return response, nil
}

func (as *authService) generateToken(user *re.User) (string, time.Time, error) {
	expiresAt := time.Now().Add(24 * time.Hour)

	claims := jwt.MapClaims{
		"userID": user.ID.String(),
		"email":  user.Email,
		"exp":    expiresAt.Unix(),
		"iat":    time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(as.jwtSecret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sign token: %w", err)
	}

	return tokenString, expiresAt, nil
}
