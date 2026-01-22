package hash

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

type Hasher interface {
	HashPassword(password string) (string, error)
	CompareHashAndPassword(hashedPassword, requestPassword string) error
}

type hasher struct{}

func NewHasher() Hasher { return &hasher{} }

func (h *hasher) HashPassword(password string) (string, error) {
	bytePassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("hash generate: %w", err)
	}

	return string(bytePassword), nil
}

func (h *hasher) CompareHashAndPassword(hashedPassword, requestPassword string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(requestPassword))
}
